// License: OpenFaaS Community Edition (CE) EULA
// Copyright (c) 2017,2019-2024 OpenFaaS Author(s)

// Copyright (c) Alex Ellis 2017. All rights reserved\.
// Copyright 2020 OpenFaaS Author(s)

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	types "github.com/danenherdi/faas-provider/types"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	v1 "k8s.io/client-go/listers/apps/v1"
	klog "k8s.io/klog"

	"github.com/openfaas/faas-netes/pkg/k8s"
)

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(defaultNamespace string, deploymentLister v1.DeploymentLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		if lookupNamespace == "kube-system" {
			http.Error(w, "unable to list within the kube-system namespace", http.StatusUnauthorized)
			return
		}

		functions, err := getServiceList(lookupNamespace, deploymentLister)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, err := json.Marshal(functions)
		if err != nil {
			klog.Errorf("Failed to marshal functions: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal functions"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

func getServiceList(functionNamespace string, deploymentLister v1.DeploymentLister) ([]types.FunctionStatus, error) {
	functions := []types.FunctionStatus{}

	sel := labels.NewSelector()
	req, err := labels.NewRequirement("faas_function", selection.Exists, []string{})
	if err != nil {
		return functions, err
	}
	onlyFunctions := sel.Add(*req)

	res, err := deploymentLister.Deployments(functionNamespace).List(onlyFunctions)

	if err != nil {
		return nil, err
	}

	for _, item := range res {
		if item != nil {
			function := k8s.AsFunctionStatus(*item)
			if function != nil {
				functions = append(functions, *function)
			}
		}
	}

	return functions, nil
}

// MakeFlowReader handler for reading workflows deployed in the cluster as deployments.
func MakeFlowReader(defaultNamespace string, deploymentLister v1.DeploymentLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
