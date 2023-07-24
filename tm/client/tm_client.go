package main

import (
	"DTM/dtm_grpc"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := ":8083"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] failed %s", addr, err))
	}
	defer conn.Close()

	// initialize client
	client := dtm_grpc.NewTransactionManagerClient(conn)
	res, err := client.Transaction(context.Background(), &dtm_grpc.TransactionRequest{
		Id1: 1,
		Id2: 2,
		Id3: 3,
	})
	fmt.Println(res, err)
}
