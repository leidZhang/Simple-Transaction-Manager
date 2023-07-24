package main

import (
	"DTM/dtm_grpc"
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
)

func dial(addr string) (conn *grpc.ClientConn, err error) {
	// create a connection to the transaction manager
	conn, err = grpc.Dial(addr, grpc.WithInsecure())

	return conn, err
}

func sendTransactionReq(client dtm_grpc.TransactionManagerClient, idA int32, idB int32, idC int32) (res *dtm_grpc.TransactionResponse, err error) {
	// create a transaction request
	req := &dtm_grpc.TransactionRequest{
		Id1: idA,
		Id2: idB,
		Id3: idC,
	}
	// send the transaction request and get response
	res, err = client.Transaction(context.Background(), req)

	return res, err
}

func expectTransactionRes(msg string, status bool) (res *dtm_grpc.TransactionResponse) {
	// create a expected response
	res = &dtm_grpc.TransactionResponse{
		Status:  status,
		Message: msg,
	}

	return res
}

func compareRes(res1 *dtm_grpc.TransactionResponse, res2 *dtm_grpc.TransactionResponse) (eq bool) {
	// get message and status
	getStatus := res1.Status
	getMsg := res1.Message
	expectStatus := res2.Status
	expectMsg := res2.Message
	// compare response
	eq = (getStatus == expectStatus) && (getMsg == expectMsg)

	return eq
}

type testCase struct {
	idA, idB, idC int32
	res           *dtm_grpc.TransactionResponse
}

func TestTm(t *testing.T) {
	// define some test cases
	tests := []testCase{
		{1, 1, 1, expectTransactionRes("Transaction success", true)},   // commit success
		{-1, 1, 1, expectTransactionRes("Transaction failed", false)},  // service A prepare failed
		{1, -1, 1, expectTransactionRes("Transaction failed", false)},  // service B prepare failed
		{1, 1, -1, expectTransactionRes("Transaction failed", false)},  // service C prepare failed
		{0, 0, 0, expectTransactionRes("Transaction failed", false)},   // all services prepare failed
		{11, 1, 1, expectTransactionRes("Transaction failed", false)},  // service A commit failed
		{1, 11, 1, expectTransactionRes("Transaction failed", false)},  // service B commit failed
		{1, 1, 11, expectTransactionRes("Transaction failed", false)},  // service C commit failed
		{0, 1, 1, expectTransactionRes("Transaction failed", false)},   // service A rollback failed
		{1, 0, 1, expectTransactionRes("Transaction failed", false)},   // service B rollback failed
		{1, 1, 0, expectTransactionRes("Transaction failed", false)},   // service C rollback failed
		{0, -1, -1, expectTransactionRes("Transaction failed", false)}, // service A prepare failed and service B and C rollback failed
		{-1, 0, -1, expectTransactionRes("Transaction failed", false)}, // service B prepare failed and service A and C rollback failed
		{-1, -1, 0, expectTransactionRes("Transaction failed", false)}, // service C prepare failed and service A and B rollback failed
	}
	// dial 8083
	conn, err := dial("localhost:8083")
	if err != nil {
		t.Fatalf(fmt.Sprintf("grpc connect addr [%s] failed %s", "localhost:8083", err))
	}
	client := dtm_grpc.NewTransactionManagerClient(conn)

	for _, test := range tests {
		getRes, err := sendTransactionReq(client, test.idA, test.idB, test.idC)
		expectRes := test.res
		if err != nil {
			t.Fatalf(err.Error(), "Error happened in transaction manager")
		}

		eq := compareRes(getRes, expectRes)
		if !eq {
			t.Fatalf("send %v, %v, %v, get: %v, but expected: %v", test.idA, test.idB, test.idC, getRes, expectRes)
		}
	}

	defer conn.Close()
}
