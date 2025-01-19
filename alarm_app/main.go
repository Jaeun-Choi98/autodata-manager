package main

import (
	"bufio"
	"cju/service/grpc_client"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/*
	이후 추가적인 기능을 개발할 때 사용.
*/

var (
	reader = bufio.NewReader(os.Stdin)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	guideStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5733"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	emailStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

	guide := `
Commands:
	subscribe <topicName>                                  - Subscribe Topic
	exit                                                   - Exit the program
`
	fmt.Println(guideStyle.Render(guide))

	for {
		fmt.Print(guideStyle.Render("\n> "))
		cmd := nextline()

		if len(cmd) == 0 {
			fmt.Println(errorStyle.Render("Invalid command. Please try again."))
			continue
		}

		switch cmd[0] {
		case "subscribe":
			handleSubscribe(ctx, cmd, emailStyle, errorStyle)
		case "exit":
			fmt.Println(successStyle.Render("Exiting the program. Goodbye!"))
			return
		}
	}
}

func handleSubscribe(ctx context.Context, cmd []string, emailStyle, errorStyle lipgloss.Style) {
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	msgSub := make(chan map[string]interface{}, 10)
	go func() {
		//defer wg.Done()
		err := grpc_client.SubscribeToMOM(ctx, cmd[1], msgSub)
		if err != nil {
			log.Println(errorStyle.Render("failed to subscribe"))
			return
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
				fmt.Println(emailStyle.Render(email))
			}
		}
	}()
	//wg.Wait()
}

func createEmail(msg map[string]interface{}) string {
	emailFormat := fmt.Sprintf(
		"Subject: %s\nFrom: sender\nTo: recipient\n%s",
		msg["Topic"],
		msg["Body"],
	)
	return emailFormat
}

func nextline() []string {
	line, _ := reader.ReadString('\n')
	return stringTokenizer(strings.TrimSpace(line))
}

func stringTokenizer(s string) []string {
	return strings.Split(s, " ")
}
