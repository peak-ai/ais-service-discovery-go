package types

// Service -
type Service struct {
	Name string
	Addr string
	Type string
}

// Signature -
type Signature struct {
	Namespace string
	Service   string
	Handler   string
}

// ServiceRequest -
type ServiceRequest struct {
	Service *Service
	Request
}

// Request -
type Request struct {
	Body []byte
}

// Response -
type Response struct {
	Body  []byte
	Error error
}

// Options type
type Options map[string]interface{}
