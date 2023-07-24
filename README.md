# Simple-Transaction-Manager
## Overview 

## Design  
The system consists of 3 service server (service A, service B, and service C), a transaction manager server, and a transaction manager client. The service server is responsible for handling specific business logic, while the transaction manager server is responsible for coordinating and processing distributed transactions. 

## Transaction processing Process 
1. Transaction client sends transaction request to the transaction manager server 
2. Transaction manager server receives transaction requests and sends prepare requests to service A, B, and C 
3. Each service performs corresponding preparation operations after receiving prepare request and returns a prepare response to the transaction manager server. 
4. Transaction manager server will collect prepare responses of all services and determines whether to proceed to the commit phase or the rollback phase based on the responses 
5. a. If all services return true prepare response, the transaction manager server will send commit request to all services, requesting them to commit transaction. <br> b. If any service return a false prepare response, the transaction manager server will send rollback request to all services, requesting them to roll back the transaction.
6. Each service execute corresponding commit or rollback operations based on requests send by the transaction manager server, and then returns a commit response or rollback response to the transaction manager server. 
7. The transaction manager server determines the final transaction status based on the received response and returns the result to the client. 

## Failure Scenarios

## Technical Debt

## Future Work

## Getting Started

