package common

// Endpoint represents an endpoint in the communication.
type Endpoint struct {
	IP      string
	Port    uint16
	Name    string
	Cmdline string
	Method  string
	Resource string
	Count int
	//Proc    string
}
