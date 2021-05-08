package impl

type ResponsePayload struct {
	RspObj      interface{}
	Exception   error
	Attachments map[string]interface{}
}
