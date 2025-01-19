package grpc_client

import (
	pb "cju/proto/v1/broker"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

func SubscribeToMOM(ctx context.Context, topic string, msgToApp chan map[string]interface{}) error {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("f	ailed to connect to MOM server: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewBrokerServiceClient(conn)

	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{Topic: topic})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
		return err
	}

	fmt.Printf("subscribed to topic: %s\n", topic)

	msgChan := make(chan *pb.SubscribeResponse, 10)
	errChan := make(chan error, 10)
	defer func() {
		close(msgChan)
		close(errChan)
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("subscription to topic '%s' canceled.\n", topic)
			return ctx.Err()
		case msg := <-msgChan:
			//fmt.Printf("received message: topic=%s, body=%s\n", msg.Topic, msg.Body)
			msgToApp <- map[string]interface{}{"Topic": msg.Topic, "Body": msg.Body}
		case err := <-errChan:
			log.Fatalf("error receiving message: %v", err)
			return err
		}
	}
}

func PublishToMOM(topic, msg string) error {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to MOM server: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewBrokerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Publish(ctx, &pb.PublishRequest{
		Topic: topic,
		Body:  msg,
	})
	if err != nil {
		log.Fatalf("failed to publish message: %v", err)
		return err
	}

	fmt.Printf("publish response: %s\n", resp.Status)
	return nil
}
