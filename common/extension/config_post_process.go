package extension

import . "github.com/laz/dubbo-go/config/interfaces"

var (
	processors = make(map[string]ConfigPostProcessor)
)

// SetConfigPostProcessor registers a ConfigPostProcessor with the given name.
func SetConfigPostProcessor(name string, processor ConfigPostProcessor) {
	processors[name] = processor
}

// GetConfigPostProcessor finds a ConfigPostProcessor by name.
func GetConfigPostProcessor(name string) ConfigPostProcessor {
	return processors[name]
}

// GetConfigPostProcessors returns all registered instances of ConfigPostProcessor.
func GetConfigPostProcessors() []ConfigPostProcessor {
	ret := make([]ConfigPostProcessor, 0, len(processors))
	for _, v := range processors {
		ret = append(ret, v)
	}
	return ret
}
