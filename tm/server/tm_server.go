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

func resMsg(feedback string, status bool) (res *dtm_grpc.TransactionResponse, err error) {
	return &dtm_grpc.TransactionResponse{
		Status:  status,
		Message: feedback,
	}, nil
}

func getPrepareRes(id int32, addr string) (conn *grpc.ClientConn, client dtm_grpc.TransactionManagerClient, prepareRes *dtm_grpc.PrepareResponse, err error) {
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

func closeConn(connA, connB, connC *grpc.ClientConn) {
	defer connA.Close()
	defer connB.Close()
	defer connC.Close()
}

func handlePartialSuccess(idA, idB, idC int32, clientA, clientB, clientC dtm_grpc.TransactionManagerClient, connA, connB, connC *grpc.ClientConn) (res *dtm_grpc.TransactionResponse, err error) {
	// send rollback request
	rollbackResA, errA := sendRollbackReq(idA, clientA)
	rollbackResB, errB := sendRollbackReq(idB, clientB)
	rollbackResC, errC := sendRollbackReq(idB, clientC)
	closeConn(connA, connB, connC)
	// check whether the services are reachable
	if errA == nil && errB == nil && errC == nil {
		fmt.Println("Service A rollback", rollbackResA)
		fmt.Println("Service B rollback", rollbackResB)
		fmt.Println("Service C rollback", rollbackResC)
	} else {
		return resMsg("Uable to reach some service at rollback stage", false)
	}
	// return false transaction response
	return resMsg("Transaction failed", false)
}

func (TransactionManager) Transaction(ctx context.Context, req *dtm_grpc.TransactionRequest) (res *dtm_grpc.TransactionResponse, err error) {
	// defer function to handle panics
	defer func() {
		if r := recover(); r != nil {
			// return the a failure response
			res, err = resMsg(fmt.Sprint(r), false)
		}
	}()

	// get transaction id
	fmt.Println(req)
	fmt.Println("Transaction request received")
	idA, idB, idC := req.Id1, req.Id2, req.Id3
	addrA, addrB, addrC := "localhost:8080", "localhost:8081", "localhost:8082"

	connA, clientA, prepareResA, errA := getPrepareRes(idA, addrA)
	connB, clientB, prepareResB, errB := getPrepareRes(idB, addrB)
	connC, clientC, prepareResC, errC := getPrepareRes(idC, addrC)
	// check whether the services are reachable
	if errA == nil && errB == nil && errC == nil {
		fmt.Println("Service A prepare", prepareResA)
		fmt.Println("Service B prepare", prepareResB)
		fmt.Println("Service C prepare", prepareResC)
	} else {
		return resMsg("Unable to reach some services at prepare stage", false)
	}

	// begin prepare
	statusA := prepareResA.Status
	statusB := prepareResB.Status
	statusC := prepareResC.Status

	// mock crash after prepare stage
	/* if idA == 3 && idB == 3 && idC == 3 {
		panic("Transaction manager internal error")
	} */

	// prepare failed, proceed to rollback
	if !statusA || !statusB || !statusC {
		return handlePartialSuccess(idA, idB, idC, clientA, clientB, clientC, connA, connB, connC)
	}
	// prepare success, proceed to commit
	commitResA, errA := sendCommitReq(idA, clientA)
	commitResB, errB := sendCommitReq(idB, clientB)
	commitResC, errC := sendCommitReq(idC, clientC)
	// check whether the services are reachable
	if errA == nil && errB == nil && errC == nil {
		fmt.Println("Service A commit", commitResA)
		fmt.Println("Service B commit", commitResB)
		fmt.Println("Service C commit", commitResC)
	} else {
		return resMsg("Unable to reach some services in commit stage", false)
	}

	// begin commit
	statusA = commitResA.Status
	statusB = commitResB.Status
	statusC = commitResC.Status
	// handle parital commit
	if !statusA || !statusB || !statusC {
		return handlePartialSuccess(idA, idB, idC, clientA, clientB, clientC, connA, connB, connC)
	}
	// transaction complete
	closeConn(connA, connB, connC)
	return resMsg("Transaction success", true)

	return res, err // in case panic happens
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
