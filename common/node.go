package common

// Node use for process dubbo node
type Node interface {
	GetUrl() *URL
	IsAvailable() bool
	Destroy()
}
