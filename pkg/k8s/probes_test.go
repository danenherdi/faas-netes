// License: OpenFaaS Community Edition (CE) EULA
// Copyright (c) 2017,2019-2024 OpenFaaS Author(s)

package k8s

import (
	"testing"

	types "github.com/danenherdi/faas-provider/types"
)

func Test_makeProbes_useExec(t *testing.T) {
	f := mockFactory()

	request := types.FunctionDeployment{
		Service:                "testfunc",
		ReadOnlyRootFilesystem: false,
	}

	probes, err := f.MakeProbes(request)
	if err != nil {
		t.Fatal(err)
	}

	if probes.Readiness.Exec == nil {
		t.Errorf("Readiness probe should have had exec handler")
		t.Fail()
	}
	if probes.Liveness.Exec == nil {
		t.Errorf("Liveness probe should have had exec handler")
		t.Fail()
	}
}

func Test_makeProbes_useHTTPProbe(t *testing.T) {
	f := mockFactory()
	f.Config.HTTPProbe = true

	request := types.FunctionDeployment{
		Service:                "testfunc",
		ReadOnlyRootFilesystem: false,
	}

	probes, err := f.MakeProbes(request)
	if err != nil {
		t.Fatal(err)
	}

	if probes.Readiness.HTTPGet == nil {
		t.Errorf("Readiness probe should have had HTTPGet handler")
		t.Fail()
	}
	if probes.Liveness.HTTPGet == nil {
		t.Errorf("Liveness probe should have had HTTPGet handler")
		t.Fail()
	}
}

func Test_makeProbes_useCustomDurationHTTPProbe(t *testing.T) {
	f := mockFactory()
	f.Config.HTTPProbe = true
	f.Config.LivenessProbe = &ProbeConfig{
		PeriodSeconds:       1,
		TimeoutSeconds:      3,
		InitialDelaySeconds: 0,
	}
	f.Config.ReadinessProbe = &ProbeConfig{
		PeriodSeconds:       1,
		TimeoutSeconds:      3,
		InitialDelaySeconds: 0,
	}

	customDelay := "0"

	request := types.FunctionDeployment{
		Service:                "testfunc",
		ReadOnlyRootFilesystem: false,
	}

	probes, err := f.MakeProbes(request)
	if err != nil {
		t.Fatal(err)
	}

	if probes.Readiness.HTTPGet == nil {
		t.Errorf("Readiness probe should have had HTTPGet handler")
		t.Fail()
	}
	if probes.Readiness.InitialDelaySeconds != 0 {
		t.Errorf("Readiness probe should have initial delay seconds set to %s", customDelay)
		t.Fail()
	}

	if probes.Liveness.HTTPGet == nil {
		t.Errorf("Liveness probe should have had HTTPGet handler")
		t.Fail()
	}
	if probes.Liveness.InitialDelaySeconds != 0 {
		t.Errorf("Readiness probe should have had HTTPGet handler set to %s", customDelay)
		t.Fail()
	}
}
