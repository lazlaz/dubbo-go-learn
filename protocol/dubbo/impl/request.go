package impl

type RequestPayload struct {
	Params      interface{}
	Attachments map[string]interface{}
}

func NewRequestPayload(args interface{}, atta map[string]interface{}) *RequestPayload {
	if atta == nil {
		atta = make(map[string]interface{})
	}
	return &RequestPayload{
		Params:      args,
		Attachments: atta,
	}
}
