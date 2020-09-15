package types

// Service represents a service endpoint as returned from a Locator
type Service struct {
	Name string
	Addr string
	Type string
}

// Signature represents a call to a service within a namespace
type Signature struct {
	Namespace string
	Service   string
	Instance  string
}

// ServiceRequest represents a call to a service
type ServiceRequest struct {
	Service *Service
	Request
}

// Request represents data to be sent to a service
type Request struct {
	Body []byte
}

// Response represents a response from a service
type Response struct {
	Body  []byte
	Error error
}

// Options is a generic key-value map that can be passed to services
type Options map[string]interface{}
