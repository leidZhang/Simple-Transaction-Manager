# Simple-Transaction-Manager
## Overview 
This project, Simple Transaction Manager, is intended to get familiar with distributed transactions and two-phase commit protocol. The transaction manager will make sure all the services will commit or rollback. <br>

I used Golang to implement the transaction manager and the relavant services, and utilized gRPC and TCP protocol for the communications between the transaction manager server and service servers.

## Design  
The system consists of 3 service server (service A, service B, and service C), a transaction manager server, and a transaction manager client. The service servers are responsible for handling specific business logic, while the transaction manager server is responsible for coordinating and processing distributed transactions. 

## Transaction Processing
1. Transaction client sends transaction request to the transaction manager server 
2. Transaction manager server receives transaction requests and sends prepare requests to service A, B, and C 
3. Each service performs corresponding preparation operations after receiving prepare request and returns a prepare response to the transaction manager server. 
4. Transaction manager server will collect prepare responses of all services and determines whether to proceed to the commit phase or the rollback phase based on the responses 
5. a. If all services return true prepare response, the transaction manager server will send commit request to all services, requesting them to commit transaction. <br> b. If any service return a false prepare response, the transaction manager server will send rollback request to all services, requesting them to roll back the transaction.
6. Each service execute corresponding commit or rollback operations based on requests send by the transaction manager server, and then returns a commit response or rollback response to the transaction manager server. 
7. The transaction manager server determines the final transaction status based on the received response and returns the result to the client. 

## Failure Scenarios
### Possible Scenarios
1. Service unreachable: In the prepare stage, the transaction manager sends a prepare request to the services, but one (or more) service is unreachable due to network or service failures, and then the transaction manager cannot receive the corresponding response. 
2. Service timeout: During the Prepare phase, some services may fail to respond in time due to network latency or high load.
3. Partial commit or rollback: During the commit or rollback phase, a partial success may occur, that is some services successfully complete commit or rollback operations while others failed to complete. 

Some of the possible scenarios is simulated in the mock services. 

### Countermeasures: 
1. Time out handling: For prepare requests and responses, set appropriate timeout times. 
2. Error rollback: If there is error or partial success, the transaction manager will send rollback requests to all services to ensure the transaction consistency. 
3. Retry: In cases of unstable network or unreachable services, retry can be performed to ensure successful sending or receiving of requests.
## Technical Debt
### Lack of fault tolerance mechanism 
In the current design, the transaction manager does not have sufficient fault-tolerant mechanisms to cope with failures or unreachable services. If one or more services become unreachable, the transaction manager may not be able to handle the failure correclty. 
### Insecure Connection 
Currently, the communication between servers is established using `grpc.Dial` with `insecure.NewCredentials()`. This results in an insecure connection without proper transport security, meaning data is transmitted in plaintext over the network. 
### Code Snipe Repetition 
The code contains repetitive error handling code snippets of the form `if err != nil { return errorMsg(err.Error()) }`. This repetitive pattern can lead to code duplication, making the code harder to maintain and prone to potential bugs. (improved) 

## Future Work
### Improve Techincal Debt
1. Improve Error Handling: Refactor the code to use custom error types or a central error handling function that handles errors consistently throughout the project. This will reduce code duplication and ensure a uniform error handling strategy.
2. Improve Connection Security: To address this technical debt, the system should be updated to use secure transport credentials to encrypt communication between servers. 
3. Improve Fault Tolerance Mechanism: Implement the following fault-tolerant mechanism: Time out handling, Error rollback, Retry 
### Extension 
1. Desgin a database for the transaction manager to store transaction status and logs, making transaction recovery and monitoring more convenient. 
2. Add a selection function to allow the transaction manager to know what kind of services is needed and send prepare requests to the corresponding services. 
3. Implement input function for the transaction manager client so that users can send transacation requests based on their requirements. 

## Getting Started
To install and run this project, you need to have the following requirements:
- golang: You need to have golang installed in your machine. For installation instructions, see <a href="https://go.dev/doc/install">Go’s Getting Started guide</a>.
- grpc: You need to have grpc installed in your machine. For installation instructions, see <a href="https://grpc.io/docs/languages/go/quickstart/">gRPC Quick Start</a>.
- protoc: You need to have the protocol buffer compiler, protoc, version 3 or higher, installed in your machine. For installation instructions, see <a href="https://grpc.io/docs/protoc-installation/">Protocol Buffer Compiler Installation</a>.
- protoc-gen-go and protoc-gen-go-grpc: You need to have the protocol compiler plugins for Go installed in your machine. For installation instructions, see <a href="https://grpc.io/docs/languages/go/basics/#generating-client-and-server-code">gRPC Basics tutorial</a>.

To clone this repository, run `git clone https://github.com/leidZhang/Simple-Transaction-Manager.git`

After the installation of the requirements and clone the repository, you can run the following command to run the server: 
- To run the transaction manager server, run `go run tm/server/tmServer.go`
- To run the serviceA server, run `go run serviceA/serverA.go`
- To run the serviceB server, run `go run serviceB/serverB.go`
- To run the serviceC server, run `go run serviceC/serverC.go`
- To send transaction request to the transaction manager server, you have to run the transaction manager client with `go run tm/client/tmClient.go`

## License
This project is licensed under the Apache-2.0 license - see the LICENSE file for details.
