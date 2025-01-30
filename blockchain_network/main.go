package main

import (
	"cju/blockchain_network/bcnet"
	pb "cju/proto/v1/bcnet"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type BCNServer struct {
	pb.UnimplementedBlockChainNetworkServiceServer
	blockChainNetwork *bcnet.BlockChainNetwork
}

func (bcn *BCNServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	bcnetReq := &bcnet.Request{Token: req.Token, Cmd: req.Cmd, Consortium: req.Consortium}
	switch req.Cmd {
	case "exit":
		if err := bcn.blockChainNetwork.ValidateRequest(bcnetReq); err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true}, nil
	case "init":
		repstr, err := bcn.blockChainNetwork.InitPeer(bcnetReq)
		if err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true, Blockchain: repstr}, nil
	case "participate":
		if err := bcn.blockChainNetwork.ValidateRequest(bcnetReq); err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true}, nil
	case "make":
		if err := bcn.blockChainNetwork.ValidateRequest(bcnetReq); err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true}, nil
	case "get":
		bc, err := bcn.blockChainNetwork.GetBlockChain(bcnetReq)
		if err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true, Blockchain: bc.BlocksStr}, nil
	case "add":
		if err := bcn.blockChainNetwork.AddBlockChain(bcnetReq, req.Transaction); err != nil {
			return &pb.MessageResponse{Success: false}, err
		}
		return &pb.MessageResponse{Success: true}, nil
	default:
		return &pb.MessageResponse{Success: false}, fmt.Errorf("bad request")
	}
}

func main() {

	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = godotenv.Load("../.env")
	if err != nil {
		log.Println(err)
	}
	bcn := bcnet.NewBlockChainNetwork()
	grpcServer := grpc.NewServer()
	pb.RegisterBlockChainNetworkServiceServer(grpcServer, &BCNServer{blockChainNetwork: bcn})

	fmt.Println("blockchain network server is running on port 9091...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
