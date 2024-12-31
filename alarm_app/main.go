package main

import (
	"cju/service/grpc_client"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

/*
	이후 추가적인 기능을 개발할 때 사용.
*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	defer cancel()
	wg.Add(1)
	msgSub := make(chan map[string]interface{}, 10)
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
				// sendEamil 함수 구현현
				email := createEmail(msg)
				fmt.Println(style.Render(email))
			}
		}
	}()
	wg.Wait()
}
func createEmail(msg map[string]interface{}) string {

	emailFormat := fmt.Sprintf(
		"Subject: %s\nFrom: sender\nTo: recipient\n%s",
		msg["Topic"],
		msg["Body"],
	)

	return emailFormat
}
