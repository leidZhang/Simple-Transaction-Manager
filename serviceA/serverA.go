package main

import (
	"DTM/dtm_grpc"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type TransactionManager struct {
}

func (TransactionManager) Transaction(ctx context.Context, req *dtm_grpc.TransactionRequest) (res *dtm_grpc.TransactionResponse, err error) {
	// will not receive transaction request
	fmt.Println(req)

	return &dtm_grpc.TransactionResponse{
		Status:  false,
		Message: "Unable to call",
	}, nil
}

func (TransactionManager) Prepare(ctx context.Context, req *dtm_grpc.PrepareRequest) (res *dtm_grpc.PrepareResponse, err error) {
	fmt.Println(req, "Preparing A")
	id := req.Id

	if id > 0 {
		return &dtm_grpc.PrepareResponse{
			Status: true,
		}, nil
	}

	return &dtm_grpc.PrepareResponse{
		Status: false,
	}, nil
}

func (TransactionManager) Commit(ctx context.Context, req *dtm_grpc.CommitRequest) (res *dtm_grpc.CommitResponse, err error) {
	// placeholder response
	fmt.Println(req, "Commiting A")
	id := req.Id

	if id > 10 {
		return &dtm_grpc.CommitResponse{
			Status: false,
		}, nil
	}

	return &dtm_grpc.CommitResponse{
		Status: true,
	}, nil
}

func (TransactionManager) Rollback(ctx context.Context, req *dtm_grpc.RollbackRequest) (res *dtm_grpc.RollbackResponse, err error) {
	// placeholder response
	fmt.Println(req, "Rollbacking A")
	id := req.Id

	if id == 0 {
		return &dtm_grpc.RollbackResponse{
			Status: false,
		}, nil
	}

	return &dtm_grpc.RollbackResponse{
		Status: true,
	}, nil
}

func main() {
	// listen port
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		grpclog.Fatalf("Failed to listen： %v", err)
	}

	// create a grpc server instance
	server := grpc.NewServer()
	service := TransactionManager{}
	dtm_grpc.RegisterTransactionManagerServer(server, &service)
	fmt.Println("grpc server running: 8080")
	err = server.Serve(listen)
}
