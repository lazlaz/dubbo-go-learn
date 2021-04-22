package protocol

import "reflect"

// Invocation is a invocation for each remote method.
type Invocation interface {
	// MethodName gets invocation method name.
	MethodName() string
	// ParameterTypeNames gets invocation parameter type names.
	ParameterTypeNames() []string
	// ParameterTypes gets invocation parameter types.
	ParameterTypes() []reflect.Type
	// ParameterValues gets invocation parameter values.
	ParameterValues() []reflect.Value
	// Arguments gets arguments.
	Arguments() []interface{}
	// Reply gets response of request
	Reply() interface{}
	// Attachments gets all attachments
	Attachments() map[string]interface{}
	// AttachmentsByKey gets attachment by key , if nil then return default value. （It will be deprecated in the future）
	AttachmentsByKey(string, string) string
	Attachment(string) interface{}
	// Attributes refers to dubbo 2.7.6.  It is different from attachment. It is used in internal process.
	Attributes() map[string]interface{}
	// AttributeByKey gets attribute by key , if nil then return default value
	AttributeByKey(string, interface{}) interface{}
	// SetAttachments sets attribute by @key and @value.
	SetAttachments(key string, value interface{})
	// Invoker gets the invoker in current context.
	Invoker() Invoker
}
