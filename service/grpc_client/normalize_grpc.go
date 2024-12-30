package grpc_client

import (
	"cju/proto/v1/normalize"
	"context"
	"log"

	"google.golang.org/grpc"
)

func NormalizeByOpenAI(tableData string) (ret string) {
	// Establish gRPC connection to the Python server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("failed to connect to gRPC server: %v", err)
	}
	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			ret = ""
		}
	}()

	client := normalize.NewNormalizeServiceClient(conn)

	req := &normalize.TableData{
		Data: tableData,
	}

	resp, err := client.NormalizeTable(context.Background(), req)
	if err != nil {
		log.Panicf("Error calling NormalizeTable: %v", err)
	}

	return resp.Result
}
