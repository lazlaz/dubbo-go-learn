package main

func main() {
	wm := NewWorkerManager(10)

	wm.StartWorkerPool()
}
