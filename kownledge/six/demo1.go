package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
	"time"
)

func rand1() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() < 0.5
}
func Rpc(ctx context.Context, url string) error {
	result := make(chan int)
	err := make(chan error)

	go func() {
		// 进行RPC调用，并且返回是否成功，成功通过result传递成功信息，错误通过error传递错误信息

		isSuccess := rand1()
		if isSuccess {
			result <- 1
		} else {
			err <- errors.New("some error happen")
		}
	}()

	select {
	case <-ctx.Done():
		// 其他RPC调用调用失败
		fmt.Println("ctx.Done")
		return ctx.Err()
	case e := <-err:
		// 本RPC调用失败，返回错误信息
		fmt.Println(e)
		return e
	case <-result:
		// 本RPC调用成功，不返回错误信息
		fmt.Println("RPC调用成功")
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// RPC1调用
	err := Rpc(ctx, "http://rpc_1_url")
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	// RPC2调用
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_2_url")
		if err != nil {
			cancel()
		}
	}()

	// RPC3调用
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_3_url")
		if err != nil {
			cancel()
		}
	}()

	// RPC4调用
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_4_url")
		if err != nil {
			cancel()
		}
	}()

	wg.Wait()
}
