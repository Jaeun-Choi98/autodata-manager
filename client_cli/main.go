package main

import (
	"bufio"
	"cju/client_cli/client"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
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
  subscribe <url>                                    - Subscribe DDL 
  unsubscribe <url>                                  - Unsubscribe DDL
  tables    <url> <schemaName>                       - Get all tables from schema
  create    <url> <fileName> <tableName> <extension> - Create a table from a file
  delete    <url> <tableName>                        - Delete a table
  read      <url> <tableName>                        - Read a table
  export    <url> <tableName> <extension>            - Export a table
  normalize <url> <fileName> <extension>             - Normalize a table from a file
  exit                                               - Exit the program
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
			if len(cmd) != 2 {
				fmt.Println(errorStyle.Render("Usage: subscribe <url>"))
			} else {
				res, err := myClient.SubscribeDDL(cmd[1])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error subscribe [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Subscribed successfully."))
					for key, val := range res {
						fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
							resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
						))
					}
				}
			}

		case "unsubscribe":
			if len(cmd) != 2 {
				fmt.Println(errorStyle.Render("Usage: unsubscribe <url>"))
			} else {
				res, err := myClient.UnsubscribeDDL(cmd[1])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error unsubscribe [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Unsubscribed successfully."))
					for key, val := range res {
						fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
							resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
						))
					}
				}
			}
		case "tables":
			if len(cmd) != 3 {
				fmt.Println(errorStyle.Render("Usage: tables <url> <schemaName>"))
			} else {
				res, err := myClient.ReadAllTables(cmd[1], cmd[2])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error reading all tables: [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Getted all tables successfully."))
					fmt.Println(resStyle.Render(fmt.Sprintf("<%s>", cmd[2])))
					if res == nil {
						fmt.Println(resStyle.Render("Nothing tables"))
					} else {
						for _, tableName := range res.([]interface{}) {
							fmt.Println(resStyle.Render(fmt.Sprintf("%s ", tableName.(string))))
						}
					}
				}
			}

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

		case "read":
			if len(cmd) != 3 {
				fmt.Println(errorStyle.Render("Usage: read <url> <tablename>"))
			} else {
				res, err := myClient.ReadAllRecord(cmd[1], cmd[2])
				if err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("Error reading table: [%v]", err)))
				} else {
					fmt.Println(successStyle.Render("Table read successfully."))
					viewTable(res["data"].([]interface{}))
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

func viewTable(datas []interface{}) {
	var columns []string
	if len(datas) > 0 {
		for data := range datas[0].(map[string]interface{}) {
			columns = append(columns, data)
		}
	}

	var records [][]string
	for _, row := range datas {
		var record []string
		for _, col := range columns {
			val := row.(map[string]interface{})[col]
			if val == nil {
				record = append(record, "")
			} else {
				record = append(record, fmt.Sprintf("%v", val))
			}
		}
		records = append(records, record)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	for _, record := range records {
		table.Append(record)
	}
	table.Render()
}
