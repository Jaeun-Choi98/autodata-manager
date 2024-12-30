package main

import (
	"cju/mom/broker"
	pb "cju/proto/v1/broker"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBrokerServiceServer
	broker *broker.Broker
}

func (s *server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	msg := broker.Message{
		Topic: req.Topic,
		Body:  req.Body,
	}
	s.broker.Publish(req.Topic, msg)
	fmt.Printf("published message: topic=%s, body=%s\n", req.Topic, req.Body)
	return &pb.PublishResponse{Status: "Message Published"}, nil
}

func (s *server) Subscribe(req *pb.SubscribeRequest, stream pb.BrokerService_SubscribeServer) error {
	topic := req.Topic
	sub := s.broker.Subscribe(topic)
	defer func() {
		s.broker.Unsubscribe(req.Topic, sub)
	}()
	fmt.Printf("new subscriber for topic: %s\n", topic)

	for {
		select {
		case msg, ok := <-sub:
			if !ok {
				return nil
			}
			if err := stream.Send(&pb.SubscribeResponse{
				Topic: req.Topic,
				Body:  msg.Body,
			}); err != nil {
				fmt.Printf("error sending message to subscriber: %v\n", err)
				return err
			}
		case <-stream.Context().Done():
			fmt.Printf("subscriber disconnected from topic: %s\n", topic)
			return stream.Context().Err()
		}
	}
}

func main() {
	broker := broker.NewBroker()

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBrokerServiceServer(grpcServer, &server{broker: broker})

	fmt.Println("MOM server is running on port 9090...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
