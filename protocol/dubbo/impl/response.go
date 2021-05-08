package impl

type ResponsePayload struct {
	RspObj      interface{}
	Exception   error
	Attachments map[string]interface{}
}

func EnsureResponsePayload(body interface{}) *ResponsePayload {
	if res, ok := body.(*ResponsePayload); ok {
		return res
	}
	if exp, ok := body.(error); ok {
		return NewResponsePayload(nil, exp, nil)
	}
	return NewResponsePayload(body, nil, nil)
}

// NewResponse create a new ResponsePayload
func NewResponsePayload(rspObj interface{}, exception error, attachments map[string]interface{}) *ResponsePayload {
	if attachments == nil {
		attachments = make(map[string]interface{})
	}
	return &ResponsePayload{
		RspObj:      rspObj,
		Exception:   exception,
		Attachments: attachments,
	}
}
