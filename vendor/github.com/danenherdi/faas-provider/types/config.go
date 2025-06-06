package types

import (
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 10 * time.Second
	defaultMaxIdleConns = 1024
)

// FaaSHandlers provide handlers for OpenFaaS
type FaaSHandlers struct {
	// ListNamespace lists namespaces which are annotated for OpenFaaS
	ListNamespaces http.HandlerFunc

	// MutateNamespace mutates a namespace to be annotated for OpenFaaS
	// each namespace must contain an annotation of "openfaas=1"
	MutateNamespace http.HandlerFunc

	// FunctionProxy provides the function invocation proxy logic.  Use proxy.NewHandlerFunc to
	// use the standard OpenFaaS proxy implementation or provide completely custom proxy logic.
	FunctionProxy http.HandlerFunc

	// Workflow Proxy provides the workflow invocation proxy logic.  Use proxy.NewHandlerFunc to
	// use the standard OpenFaaS proxy implementation or provide completely custom proxy logic.
	Flows      http.HandlerFunc
	FlowProxy  http.HandlerFunc
	FlowReader http.HandlerFunc

	// FunctionLister lists deployed functions within a namespace
	FunctionLister http.HandlerFunc

	// DeployFunction deploys a function which doesn't exist
	DeployFunction http.HandlerFunc

	// UpdateFunction updates an existing function
	UpdateFunction http.HandlerFunc

	DeleteFunction http.HandlerFunc

	FunctionStatus http.HandlerFunc

	ScaleFunction http.HandlerFunc

	Secrets http.HandlerFunc

	// Logs provides streaming json logs of functions
	Logs http.HandlerFunc

	// Health defines the default health endpoint bound to "/healthz
	// If the handler is not set, then the "/healthz" path will not be configured
	Health http.HandlerFunc

	Info http.HandlerFunc

	Telemetry http.HandlerFunc
}

// FaaSConfig set config for HTTP handlers
type FaaSConfig struct {
	// TCPPort is the public port for the API.
	TCPPort *int
	// HTTP timeout for reading a request from clients.
	ReadTimeout time.Duration
	// HTTP timeout for writing a response from functions.
	WriteTimeout time.Duration
	// EnableHealth enables/disables the default health endpoint bound to "/healthz".
	//
	// Deprecated: basic auth is enabled automatcally by setting the HealthHandler in the FaaSHandlers
	// struct.  This value is not longer read or used.
	EnableHealth bool
	// EnableBasicAuth enforces basic auth on the API. If set, reads secrets from file-system
	// location specificed in `SecretMountPath`.
	EnableBasicAuth bool
	// SecretMountPath specifies where to read secrets from for embedded basic auth.
	SecretMountPath string
	// MaxIdleConns with a default value of 1024, can be used for tuning HTTP proxy performance.
	MaxIdleConns int
	// MaxIdleConnsPerHost with a default value of 1024, can be used for tuning HTTP proxy performance.
	MaxIdleConnsPerHost int
	// EnableCaching
	EnableCaching bool
}

// GetReadTimeout is a helper to safely return the configured ReadTimeout or the default value of 10s
func (c *FaaSConfig) GetReadTimeout() time.Duration {
	if c.ReadTimeout <= 0*time.Second {
		return defaultReadTimeout
	}
	return c.ReadTimeout
}

// GetMaxIdleConns is a helper to safely return the configured MaxIdleConns or the default value of 1024
func (c *FaaSConfig) GetMaxIdleConns() int {
	if c.MaxIdleConns < 1 {
		return defaultMaxIdleConns
	}

	return c.MaxIdleConns
}

// GetMaxIdleConns is a helper to safely return the configured MaxIdleConns or the default value which
// should then match the MaxIdleConns
func (c *FaaSConfig) GetMaxIdleConnsPerHost() int {
	if c.MaxIdleConnsPerHost < 1 {
		return c.GetMaxIdleConns()
	}

	return c.MaxIdleConnsPerHost
}
