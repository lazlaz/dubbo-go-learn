package protocol

// Result is a RPC result
type Result interface {
	// SetError sets error.
	SetError(error)
	// Error gets error.
	Error() error
	// SetResult sets invoker result.
	SetResult(interface{})
	// Result gets invoker result.
	Result() interface{}
	// SetAttachments replaces the existing attachments with the specified param.
	SetAttachments(map[string]interface{})
	// Attachments gets all attachments
	Attachments() map[string]interface{}

	// AddAttachment adds the specified map to existing attachments in this instance.
	AddAttachment(string, interface{})
	// Attachment gets attachment by key with default value.
	Attachment(string, interface{}) interface{}
}
