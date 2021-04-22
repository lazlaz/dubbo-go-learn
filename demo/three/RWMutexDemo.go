package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mutex *sync.RWMutex
	mutex = new(sync.RWMutex)
	fmt.Println("Lock the lock")
	mutex.Lock()
	fmt.Println("The lock is locked")

	channels := make([]chan int, 4)
	for i := 0; i < 4; i++ {
		channels[i] = make(chan int)
		go func(i int, c chan int) {
			fmt.Println("Not read lock: ", i)
			mutex.RLock()
			fmt.Println("Read Locked: ", i)
			fmt.Println("Unlock the read lock: ", i)
			time.Sleep(time.Second)
			mutex.RUnlock()
			c <- i
		}(i, channels[i])
	}
	time.Sleep(time.Second)
	fmt.Println("Unlock the lock")
	mutex.Unlock()
	time.Sleep(time.Second)

	for _, c := range channels {
		<-c
	}
}
