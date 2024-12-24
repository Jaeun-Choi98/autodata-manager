package main

import (
	"bufio"
	"cju/client_cli/client"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/*
https://github.com/charmbracelet/bubbletea/tree/master/examples/table
콘솔로 테이블 표를 만들려면 위 url 참고.
*/

var (
	reader   = bufio.NewReader(os.Stdin)
	myClient = client.NewClient()
)

func main() {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5733"))

	guide := `
Commands:
  create <url> <filename> <tablename> <extension> - Create a table from a file
  delete <url> <tablename>                       - Delete a table
  export <url> <tablename> <extension>           - Export a table
  normalize <url> <filename> <extension>         - Normalize a table from a file
  exit                                           - Exit the program
`
	fmt.Println(style.Render(guide))

	for {
		fmt.Print(style.Render("\n> "))
		cmd := nextline()

		if len(cmd) == 0 {
			fmt.Println("Invalid command. Please try again.")
			continue
		}

		switch cmd[0] {
		case "create":
			if len(cmd) != 5 {
				fmt.Println("Usage: create <url> <filename> <tablename> <extension>")
			} else {
				err := myClient.MakeTable(cmd[1], cmd[2], cmd[3], cmd[4])
				if err != nil {
					fmt.Printf("Error creating table: %v\n", err)
				} else {
					fmt.Println("Table created successfully.")
				}
			}

		case "delete":
			if len(cmd) != 3 {
				fmt.Println("Usage: delete <url> <tablename>")
			} else {
				err := myClient.DropTable(cmd[1], cmd[2])
				if err != nil {
					fmt.Printf("Error deleting table: %v\n", err)
				} else {
					fmt.Println("Table deleted successfully.")
				}
			}

		case "export":
			if len(cmd) != 4 {
				fmt.Println("Usage: export <url> <tablename> <extension>")
			} else {
				err := myClient.ExportTable(cmd[1], cmd[2], cmd[3])
				if err != nil {
					fmt.Printf("Error exporting table: %v\n", err)
				} else {
					fmt.Println("Table exported successfully.")
				}
			}

		case "normalize":
			if len(cmd) != 4 {
				fmt.Println("Usage: normalize <url> <filename> <extension>")
			} else {
				err := myClient.NormalizeTable(cmd[1], cmd[2], cmd[3])
				if err != nil {
					fmt.Printf("Error normalizing table: %v\n", err)
				} else {
					fmt.Println("Table normalized successfully.")
				}
			}

		case "exit":
			fmt.Println("Exiting the program. Goodbye!")
			return

		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

func nextline() []string {
	line, _ := reader.ReadString('\n')
	return stringTokenizer(strings.TrimSpace(line))
}

func stringTokenizer(s string) []string {
	return strings.Split(s, " ")
}
