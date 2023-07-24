package main

import (
	"DTM/dtm_grpc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type TransactionManager struct {
}

func errorMsg(feedback string) (res *dtm_grpc.TransactionResponse, err error) {
	return &dtm_grpc.TransactionResponse{
		Status:  false,
		Message: feedback,
	}, nil
}

func getPrepareResponse(id int32, addr string) (conn *grpc.ClientConn, client dtm_grpc.TransactionManagerClient, prepareRes *dtm_grpc.PrepareResponse, err error) {
	// create grpc connection
	conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf(fmt.Sprintf("grpc connect addr [%s] failed %s", addr, err))
	}
	// initialize client
	client = dtm_grpc.NewTransactionManagerClient(conn)
	prepareReq := &dtm_grpc.PrepareRequest{
		Id: id,
	}
	prepareRes, err = client.Prepare(context.Background(), prepareReq)

	return conn, client, prepareRes, err
}

func sendRollbackReq(id int32, client dtm_grpc.TransactionManagerClient) (rollbackRes *dtm_grpc.RollbackResponse, err error) {
	rollbackReq := &dtm_grpc.RollbackRequest{
		Id: id,
	}
	rollbackRes, err = client.Rollback(context.Background(), rollbackReq)

	return rollbackRes, err
}

func sendCommitReq(id int32, client dtm_grpc.TransactionManagerClient) (commitRes *dtm_grpc.CommitResponse, err error) {
	commitReq := &dtm_grpc.CommitRequest{
		Id: id,
	}
	commitRes, err = client.Commit(context.Background(), commitReq)

	return commitRes, err
}

func closeConn(connA *grpc.ClientConn, connB *grpc.ClientConn, connC *grpc.ClientConn) {
	defer connA.Close()
	defer connB.Close()
	defer connC.Close()
}

func (TransactionManager) Transaction(ctx context.Context, req *dtm_grpc.TransactionRequest) (res *dtm_grpc.TransactionResponse, err error) {
	// get transaction id
	fmt.Println(req)
	fmt.Println("Transaction request received")
	idA, idB, idC := req.Id1, req.Id2, req.Id3
	addrA, addrB, addrC := "localhost:8080", "localhost:8081", "localhost:8082"

	connA, clientA, prepareResA, err := getPrepareResponse(idA, addrA)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service A prepare", prepareResA)

	connB, clientB, prepareResB, err := getPrepareResponse(idB, addrB)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service B prepare", prepareResB)

	connC, clientC, prepareResC, err := getPrepareResponse(idC, addrC)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service C prepare", prepareResC)

	// begin prepare
	statusA := prepareResA.Status
	statusB := prepareResB.Status
	statusC := prepareResC.Status
	// prepare failed, proceed to rollback
	if !statusA || !statusB || !statusC {
		rollbackResA, err := sendRollbackReq(idA, clientA)
		if err != nil {
			return errorMsg(err.Error())
		}
		fmt.Println("Service A rollback", rollbackResA)
		rollbackResB, err := sendRollbackReq(idB, clientB)
		if err != nil {
			return errorMsg(err.Error())
		}
		fmt.Println("Service B rollback", rollbackResB)
		rollbackResC, err := sendRollbackReq(idB, clientC)
		if err != nil {
			return errorMsg(err.Error())
		}
		fmt.Println("Service C rollback", rollbackResC)

		closeConn(connA, connB, connC)
		fmt.Println("Transaction Failed")
		return &dtm_grpc.TransactionResponse{
			Status:  false,
			Message: "Transaction Failed",
		}, nil
	}
	// prepare success, proceed to commit
	commitResA, err := sendCommitReq(idA, clientA)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service A commit", commitResA)
	commitResB, err := sendCommitReq(idB, clientB)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service B commit", commitResB)
	commitResC, err := sendCommitReq(idC, clientC)
	if err != nil {
		return errorMsg(err.Error())
	}
	fmt.Println("Service C commit", commitResC)

	// transaction complete
	closeConn(connA, connB, connC)
	fmt.Println("Transaction Success")
	return &dtm_grpc.TransactionResponse{
		Status:  true,
		Message: "Transaction Success",
	}, nil
}

func (TransactionManager) Prepare(ctx context.Context, req *dtm_grpc.PrepareRequest) (res *dtm_grpc.PrepareResponse, err error) {
	// will not receive prepare request
	fmt.Println(req)

	return &dtm_grpc.PrepareResponse{
		Status: false,
	}, nil
}

func (TransactionManager) Commit(ctx context.Context, req *dtm_grpc.CommitRequest) (res *dtm_grpc.CommitResponse, err error) {
	// will not receive commit request
	fmt.Println(req)

	return &dtm_grpc.CommitResponse{
		Status: false,
	}, nil
}

func (TransactionManager) Rollback(ctx context.Context, req *dtm_grpc.RollbackRequest) (res *dtm_grpc.RollbackResponse, err error) {
	// will not receive rollback request
	fmt.Println(req)

	return &dtm_grpc.RollbackResponse{
		Status: false,
	}, nil
}

func main() {
	// listen port
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		grpclog.Fatalf("Failed to listen： %v", err)
	}

	// create a grpc server instance
	server := grpc.NewServer()
	service := TransactionManager{}
	dtm_grpc.RegisterTransactionManagerServer(server, &service)
	fmt.Println("grpc server running: 8083")
	err = server.Serve(listen)
}
