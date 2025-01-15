package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateTableFromCSV(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(dbHost, dbPort, dbPwd, dbName)
	if err != nil {
		log.Println("NewService method err")
		return
	}
	err = service.CreateTableFromCSV("../data.csv", "testtable")
	if err != nil {
		log.Println("Service.CreateTableFromCSV err")
	}
}

func TestReadAllRecordByTableName(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(dbHost, dbPort, dbPwd, dbName)
	if err != nil {
		log.Println("NewService method err")
		return
	}
	ret, _ := service.ReadAllRecordByTableName("testtable")
	fmt.Printf("return val: %v", ret)
}

func TestListenerManager(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s", dbPwd, dbHost, dbPort, dbName)

	lm, err := NewListenManager(connStr)
	testCon, _ := pgx.Connect(context.Background(), connStr)
	if err != nil {
		t.Fatalf("Failed to create ListenerManager: %v", err)
	}
	defer func() {
		lm.Close()
		testCon.Close(context.Background())
	}()

	t.Run("StartListening", func(t *testing.T) {
		err := lm.StartListening()
		assert.NoError(t, err, "Should start listening without error")

		err = lm.StartListening()
		assert.Error(t, err, "Should not start listening again")
	})

	t.Run("StopListening", func(t *testing.T) {
		lm.StopListening()
		// 리스닝이 중지되었는지 확인
		// 이 부분은 별도로 상태를 추적하는 로직을 추가해야 할 수도 있음
		// 예시로는 lm.listening 값을 활용할 수 있음.
		// assert.False(t, lm.listening, "Should stop listening after StopListening is called")
	})

	t.Run("listenLoop", func(t *testing.T) {
		err := lm.StartListening()
		assert.NoError(t, err)

		// 알림을 보내는 부분 (실제로 데이터베이스에 트리거가 설정되어 있어야 함)
		go func() {
			time.Sleep(1 * time.Second)
			_, err := testCon.Exec(context.Background(), "NOTIFY table_events, 'Test Notification'")
			if err != nil {
				t.Errorf("Failed to send notification: %v", err)
			}
		}()

		// 알림을 받을 수 있도록 대기
		select {
		case notification := <-lm.notificationChan:
			assert.Equal(t, "Test Notification", notification, "Should receive the correct notification")
		case <-time.After(2 * time.Second):
			t.Error("Test notification was not received in time")
		}

		lm.StopListening()
	})

	t.Run("Close", func(t *testing.T) {
		err := lm.Close()
		assert.NoError(t, err, "Should close the connection without error")
	})
}

func TestBackupDatabase(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(dbHost, dbPort, dbPwd, dbName)
	if err != nil {
		log.Println("NewService method err")
		return
	}
	err = service.BackupDatabase("test")
	assert.NoError(t, err, "asdf")
}

func TestAddUserFromCSV(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(dbHost, dbPort, dbPwd, dbName)
	if err != nil {
		log.Println("NewService method err")
		return
	}
	err = service.AddUserFromCSV("../member_data.csv")
	assert.NoError(t, err, "asdf")
}
