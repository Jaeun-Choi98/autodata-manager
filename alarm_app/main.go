package main

import (
	"cju/service/grpc_client"
	"context"
	"fmt"
	"log"
	"sync"
)

/*
	이후 추가적인 기능을 개발할 때 사용.
*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	defer cancel()
	wg.Add(1)
	msgSub := make(chan map[string]interface{})
	go func() {
		defer wg.Done()
		err := grpc_client.SubscribeToMOM(ctx, "table_events", msgSub)
		if err != nil {
			log.Println("failed to subscribe")
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgSub:
				fmt.Printf("topic: %s , body: %s \n", msg["Topic"], msg["Body"])
				// sendEamil 함수 구현현
			}
		}
	}()
	wg.Wait()
}
