package main

import (
	"bufio"
	"cju/client_cli/client"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	reader   = bufio.NewReader(os.Stdin)
	myClient = client.NewClient()
)

func main() {
	// Guide 스타일
	guideStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5733"))

	// 성공 스타일
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))

	// 에러 스타일
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000"))

	// 결과 스타일
	resStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	guide := `
Commands:
  create <url> <filename> <tablename> <extension> - Create a table from a file
  delete <url> <tablename>                       - Delete a table
  export <url> <tablename> <extension>           - Export a table
  normalize <url> <filename> <extension>         - Normalize a table from a file
  exit                                           - Exit the program
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
		case "create":
			if len(cmd) != 5 {
				fmt.Println(errorStyle.Render("Usage: create <url> <filename> <tablename> <extension>"))
			} else {
				res, err := myClient.MakeTable(cmd[1], cmd[2], cmd[3], cmd[4])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating table [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Table created successfully."))
					for key, val := range res {
						fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
							resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
						))
					}
				}
			}

		case "delete":
			if len(cmd) != 3 {
				fmt.Println(errorStyle.Render("Usage: delete <url> <tablename>"))
			} else {
				res, err := myClient.DropTable(cmd[1], cmd[2])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error deleting table: [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Table deleted successfully."))
					for key, val := range res {
						fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
							resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
						))
					}
				}
			}

		case "export":
			if len(cmd) != 4 {
				fmt.Println(errorStyle.Render("Usage: export <url> <tablename> <extension>"))
			} else {
				err := myClient.ExportTable(cmd[1], cmd[2], cmd[3])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error exporting table: [%v]", err)))
				} else {
					fmt.Println(successStyle.Render(fmt.Sprintf("Table exported successfully. (%s.%s)", cmd[2], cmd[3])))
				}
			}

		case "normalize":
			if len(cmd) != 4 {
				fmt.Println(errorStyle.Render("Usage: normalize <url> <filename> <extension>"))
			} else {
				res, err := myClient.NormalizeTable(cmd[1], cmd[2], cmd[3])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error normalizing table: [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Table normalized successfully."))
					for key, val := range res {
						fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
							resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
						))
					}
				}
			}

		case "exit":
			fmt.Println(successStyle.Render("Exiting the program. Goodbye!"))
			return

		default:
			fmt.Println(errorStyle.Render("Invalid command. Please try again."))
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
