// License: OpenFaaS Community Edition (CE) EULA
// Copyright (c) 2017,2019-2024 OpenFaaS Author(s)

// Copyright (c) Alex Ellis 2017. All rights reserved.
// Copyright 2020 OpenFaaS Author(s)

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	types "github.com/danenherdi/faas-provider/types"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-netes/pkg/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/client-go/listers/apps/v1"
	klog "k8s.io/klog"
)

// MaxReplicas licensed for OpenFaaS CE is 5/5
// a license for OpenFaaS Standard is required to increase this limit.
const MaxReplicas = 5

// MaxFunctions licensed for OpenFaaS CE is 15
// a license for OpenFaaS Standard is required to increase this limit.
const MaxFunctions = 15

// MakeReplicaReader reads the amount of replicas for a deployment
func MakeReplicaReader(defaultNamespace string, lister v1.DeploymentLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		if lookupNamespace != defaultNamespace {
			http.Error(w, fmt.Sprintf("namespace must be: %s", defaultNamespace), http.StatusBadRequest)
			return
		}

		s := time.Now()

		function, err := getService(lookupNamespace, functionName, lister)
		if err != nil {
			log.Printf("Unable to fetch service: %s", functionName)
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("Unable to fetch service: %s", functionName), http.StatusInternalServerError)
			return
		}

		if function == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		d := time.Since(s)
		log.Printf("Replicas: %s.%s, (%d/%d) %dms", functionName, lookupNamespace, function.AvailableReplicas, function.Replicas, d.Milliseconds())

		functionBytes, err := json.Marshal(function)
		if err != nil {
			klog.Errorf("Failed to marshal function: %s", err.Error())
			http.Error(w, "Failed to marshal function", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

// getService returns a function/service or nil if not found
func getService(functionNamespace string, functionName string, lister v1.DeploymentLister) (*types.FunctionStatus, error) {

	item, err := lister.Deployments(functionNamespace).
		Get(functionName)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if item != nil {
		function := k8s.AsFunctionStatus(*item)
		if function != nil {
			return function, nil
		}
	}

	return nil, fmt.Errorf("function: %s not found", functionName)
}
