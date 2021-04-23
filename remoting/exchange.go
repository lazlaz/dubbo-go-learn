package remoting

// this is request for transport layer
type Request struct {
	ID int64
	// protocol version
	Version string
	// serial ID (ignore)
	SerialID byte
	// Data
	Data   interface{}
	TwoWay bool
	Event  bool
}

// this is response for transport layer
type Response struct {
	ID       int64
	Version  string
	SerialID byte
	Status   uint8
	Event    bool
	Error    error
	Result   interface{}
}
