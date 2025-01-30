package client

import (
	pb "cju/proto/v1/bcnet"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Consortium struct {
	ConsortiumExist map[string]bool
	ConsortiumSlice []string
}

func NewConsortium() *Consortium {
	return &Consortium{
		ConsortiumExist: make(map[string]bool),
		ConsortiumSlice: make([]string, 0),
	}
}

func (c *Consortium) Add(val string) {
	if _, exists := c.ConsortiumExist[val]; !exists {
		c.ConsortiumExist[val] = true
		c.ConsortiumSlice = append(c.ConsortiumSlice, val)
	}
}

func (c *Consortium) Exists(val string) bool {
	_, exists := c.ConsortiumExist[val]
	return exists
}

func (c *Consortium) Remove(val string) {
	if _, exists := c.ConsortiumExist[val]; exists {
		delete(c.ConsortiumExist, val)
		newSlice := []string{}
		for _, v := range c.ConsortiumSlice {
			if v != val {
				newSlice = append(newSlice, v)
			}
		}
		c.ConsortiumSlice = newSlice
	}
}

type BcClient struct {
	con         *grpc.ClientConn
	client      pb.BlockChainNetworkServiceClient
	Consortiums *Consortium
}

func NewBlockChainClient() (*BcClient, error) {
	con, err := grpc.NewClient("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &BcClient{
		con:         con,
		client:      pb.NewBlockChainNetworkServiceClient(con),
		Consortiums: NewConsortium(),
	}, nil
}

func (bc *BcClient) Close() {
	bc.con.Close()
}

func (bc *BcClient) Do(req *pb.MessageRequest) (*pb.MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := bc.client.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (bc *BcClient) InitPeer(consortiums []string) {
	bc.Consortiums.ConsortiumExist = make(map[string]bool)
	bc.Consortiums.ConsortiumSlice = make([]string, 0)
	for _, consortium := range consortiums {
		bc.Consortiums.ConsortiumExist[consortium] = true
		bc.Consortiums.ConsortiumSlice = append(bc.Consortiums.ConsortiumSlice, consortium)
	}
}

// 클라이언트의 모든 컨소시엄에 트랜잭션 데이터를 보냄
func (bc *BcClient) SendAllConsortiumsTransactions(token string, data []string) error {
	for _, consortium := range bc.Consortiums.ConsortiumSlice {
		pbRes, err := bc.Do(&pb.MessageRequest{Token: token, Cmd: "add", Transaction: data, Consortium: consortium})
		if err != nil {
			return err
		}
		if !pbRes.Success {
			return fmt.Errorf("failed to send transactions")
		}
	}
	return nil
}

// 특정 컨소시엄에만 트랜잭션 데이터를 보낼 때
func (bc *BcClient) SendTransactions(token, consortium string, data []string) error {
	pbRes, err := bc.Do(&pb.MessageRequest{Token: token, Cmd: "add", Transaction: data, Consortium: consortium})
	if err != nil {
		return err
	}
	if !pbRes.Success {
		return fmt.Errorf("failed to send transactions")
	}
	return nil
}
