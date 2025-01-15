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
	// 스타일 정의
	guideStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5733"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
	resStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))

	guide := `
Commands:
	cron <option> [start | stop | remove <jobId> | jobs]   - Manage cron jobs
	backup <dbName>                                        - Backup a database
	cronbackup <dbName> <cronQuery>                        - Set up cron backup with a query
	listen                                                 - Listen to DDL changes
	unlisten                                               - Unlisten from DDL changes
	tables <schemaName>                                    - List all tables in a schema
	create <fileName> <tableName> <extension>              - Create a table from a file
	delete <tableName>                                     - Delete a table
	read <tableName>                                       - Read data from a table
	export <tableName> <extension>                         - Export a table to a file
	normalize <fileName> <extension>                       - Normalize a table from a file
	schema create <schemaName>                             - Create a schema
	schema delete <schemaName> <option> [-f]               - Delete a schema
	schema list                                            - List all schemas in a database
	register <fileName>                                    - Register users 
	update <fileName>                                      - Update users 
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
		case "schema":
			handleSchemaCmd(cmd, successStyle, errorStyle, resStyle)
		case "cronbackup":
			handleCronBackup(cmd, successStyle, errorStyle, resStyle)
		case "backup":
			handleBackup(cmd, successStyle, errorStyle, resStyle)
		case "cron":
			handleCronCommand(cmd, successStyle, errorStyle, resStyle)
		case "listen":
			handleListen(cmd, successStyle, errorStyle, resStyle)
		case "unlisten":
			handleUnlisten(cmd, successStyle, errorStyle, resStyle)
		case "tables":
			handleTables(cmd, successStyle, errorStyle, resStyle)
		case "create":
			handleCreate(cmd, successStyle, errorStyle, resStyle)
		case "read":
			handleRead(cmd, successStyle, errorStyle, resStyle)
		case "delete":
			handleDelete(cmd, successStyle, errorStyle, resStyle)
		case "export":
			handleExport(cmd, successStyle, errorStyle, resStyle)
		case "normalize":
			handleNormalize(cmd, successStyle, errorStyle, resStyle)
		case "register":
			handleRegister(cmd, successStyle, errorStyle, resStyle)
		case "update":
			handleUpdate(cmd, successStyle, errorStyle, resStyle)
		case "exit":
			fmt.Println(successStyle.Render("Exiting the program. Goodbye!"))
			return
		default:
			fmt.Println(errorStyle.Render("Invalid command. Please try again."))
		}
	}
}

// 각 명령어 처리 함수들

func handleRegister(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) < 1 {
		fmt.Println(errorStyle.Render("Usage: register <fileName>"))
		return
	}
	res, err := myClient.RegisterUser(cmd[1])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleUpdate(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) < 1 {
		fmt.Println(errorStyle.Render("Usage: update <fileName>"))
		return
	}
	res, err := myClient.UpdateUser(cmd[1])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleSchemaCmd(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) < 2 {
		fmt.Println(errorStyle.Render("Usage: schema <cmd> [create | delete | list]"))
		return
	}
	switch cmd[1] {
	case "create":
		res, err := myClient.MakeSchema(cmd[2])
		handleResponse(res, err, successStyle, errorStyle, resStyle)
	case "delete":
		res, err := myClient.DropSchema(cmd[2])
		handleResponse(res, err, successStyle, errorStyle, resStyle)
	case "list":
		res, err := myClient.ReadAllSchemas()
		handleResponse(res, err, successStyle, errorStyle, resStyle)
	default:
		fmt.Println(errorStyle.Render("Usage: schema <cmd> [create | delete | list]"))
		return
	}
}

func handleCronBackup(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) < 7 {
		fmt.Println(errorStyle.Render("Usage: cronbackup <dbName> <cronQuery>"))
		return
	}

	query := cmd[2:7]
	res, err := myClient.CronBackupDB(cmd[1], query)
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleBackup(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 2 {
		fmt.Println(errorStyle.Render("Usage: backup <dbName>"))
		return
	}

	res, err := myClient.BackupDB(cmd[1])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleCronCommand(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) < 2 {
		fmt.Println(errorStyle.Render("Usage: cron <option> [start | stop | remove <jobId> | jobs]"))
		return
	}

	switch cmd[1] {
	case "remove":
		if len(cmd) < 3 {
			fmt.Println(errorStyle.Render("Usage: cron remove <jobId>"))
			return
		}
		res, err := myClient.CronCommand(cmd[1], cmd[2])
		handleResponse(res, err, successStyle, errorStyle, resStyle)

	case "jobs":
		res, err := myClient.CronCommand(cmd[1], "")
		handleResponse(res, err, successStyle, errorStyle, resStyle)

	default:
		_, err := myClient.CronCommand(cmd[1], "")
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error: [%v]", err)))
		} else {
			fmt.Println(successStyle.Render("Operation completed successfully."))
		}
	}
}

func handleListen(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 1 {
		fmt.Println(errorStyle.Render("Usage: subscribe"))
		return
	}

	res, err := myClient.Listen()
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleUnlisten(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 1 {
		fmt.Println(errorStyle.Render("Usage: unsubscribe"))
		return
	}

	res, err := myClient.Unlisten()
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleTables(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 2 {
		fmt.Println(errorStyle.Render("Usage: tables <schemaName>"))
		return
	}

	res, err := myClient.ReadAllTables(cmd[1])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleCreate(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 4 {
		fmt.Println(errorStyle.Render("Usage: create <fileName> <tableName> <extension>"))
		return
	}

	res, err := myClient.MakeTable(cmd[1], cmd[2], cmd[3])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleRead(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 2 {
		fmt.Println(errorStyle.Render("Usage: read <tableName>"))
		return
	}

	res, err := myClient.ReadAllRecord(cmd[1])
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error reading table: [%v]", err)))
	} else {
		printTableResponse(res["data"].([]interface{}), successStyle, resStyle)
	}
}

func handleDelete(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 2 {
		fmt.Println(errorStyle.Render("Usage: delete <tableName>"))
		return
	}

	res, err := myClient.DropTable(cmd[1])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

func handleExport(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 3 {
		fmt.Println(errorStyle.Render("Usage: export <tableName> <extension>"))
		return
	}

	err := myClient.ExportTable(cmd[1], cmd[2])
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error exporting table: [%v]", err)))
	} else {
		fmt.Println(successStyle.Render(fmt.Sprintf("Table exported successfully. (%s.%s)", cmd[1], cmd[2])))
	}
}

func handleNormalize(cmd []string, successStyle, errorStyle, resStyle lipgloss.Style) {
	if len(cmd) != 3 {
		fmt.Println(errorStyle.Render("Usage: normalize <fileName> <extension>"))
		return
	}

	res, err := myClient.NormalizeTable(cmd[1], cmd[2])
	handleResponse(res, err, successStyle, errorStyle, resStyle)
}

// 공통 응답 처리 함수
func handleResponse(res interface{}, err error, successStyle, errorStyle, resStyle lipgloss.Style) {
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error: [%v]", err)))
	} else {
		printSuccessResponse(res, successStyle, resStyle)
	}
}

// 결과 출력 함수
func printSuccessResponse(res interface{}, successStyle, resStyle lipgloss.Style) {
	fmt.Println(successStyle.Render("Operation completed successfully."))
	for key, val := range res.(map[string]interface{}) {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
			resStyle.Render(fmt.Sprintf("%v: %v", key, val)),
		))
	}
}

// 테이블 출력 함수
func printTableResponse(datas []interface{}, successStyle, resStyle lipgloss.Style) {
	fmt.Println(successStyle.Render("Operation completed successfully."))
	if datas == nil {
		fmt.Println(resStyle.Render("No tables found"))
	} else {
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
}

func nextline() []string {
	line, _ := reader.ReadString('\n')
	return stringTokenizer(strings.TrimSpace(line))
}

func stringTokenizer(s string) []string {
	return strings.Split(s, " ")
}
