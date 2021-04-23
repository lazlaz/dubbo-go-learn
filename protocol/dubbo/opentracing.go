package dubbo

func filterContext(attachments map[string]interface{}) map[string]string {
	var traceAttchment = make(map[string]string)
	for k, v := range attachments {
		if r, ok := v.(string); ok {
			traceAttchment[k] = r
		}
	}
	return traceAttchment
}
