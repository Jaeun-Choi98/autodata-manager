package service

import (
	"cju/service/grpc_client"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
)

type ListenerManagerInterface interface {
	StartListening() error
	StopListening() error
	Close() error
}

type ListenerManager struct {
	conn             *pgx.Conn
	stopListening    chan bool
	notificationChan chan string
	wg               sync.WaitGroup
	mu               sync.Mutex
	listening        bool
}

func NewListenManager(conInfo string) (*ListenerManager, error) {
	conn, err := pgx.Connect(context.Background(), conInfo)
	if err != nil {
		log.Println("failed to connect to db(ListenerManager)")
		return nil, err
	}
	return &ListenerManager{
		conn:             conn,
		stopListening:    make(chan bool),
		notificationChan: make(chan string, 10),
	}, nil
}

func (lm *ListenerManager) StartListening() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.listening {
		log.Println("already listening")
		return fmt.Errorf("already listening")
	}

	if lm.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}

	_, err := lm.conn.Exec(context.Background(), "LISTEN table_events")
	if err != nil {
		log.Printf("failed to execute LISTEN cmd: %v", err)
		return err
	}

	lm.listening = true
	lm.wg.Add(1)
	go lm.listenLoop()

	return nil
}

func (lm *ListenerManager) listenLoop() {
	defer lm.wg.Done()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notificationWaitGroup := sync.WaitGroup{}
	notificationWaitGroup.Add(1)

	fmt.Println("Listening for table change event...")
	go func() {
		defer notificationWaitGroup.Done()
		for {
			notification, err := lm.conn.WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Notification wait canceled")
					return
				}
				log.Printf("failed to wait for notification: %v", err)
				continue
			}
			lm.notificationChan <- notification.Payload
		}
	}()

	for {
		select {
		case <-lm.stopListening:
			fmt.Println("Stop signal received. Exiting listen loop...")
			cancel()
			notificationWaitGroup.Wait()
			_, err := lm.conn.Exec(context.Background(), "UNLISTEN table_events")
			if err != nil {
				log.Printf("failed to unlisten: %v", err)
			} else {
				fmt.Println("Stopped listening to notifications.")
			}
			return
		case payload := <-lm.notificationChan:
			fmt.Printf("Received notification: %s\n", payload)
			// 이후 메시지 알람 기능 구현
			go grpc_client.PublishToMOM("table_events", payload)
		}
	}
}

func (lm *ListenerManager) StopListening() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if !lm.listening {
		log.Println("not currently listening")
		return fmt.Errorf("not currently listening")
	}

	lm.stopListening <- true
	lm.wg.Wait()
	lm.listening = false
	return nil
}

func (lm *ListenerManager) Close() error {
	err := lm.conn.Close(context.Background())
	if err != nil {
		log.Println("failed to close db(ListenManager)")
		return err
	}
	return nil
}
