# go-clean-architecture

Esta listagem precisa ser feita com:

- Endpoint REST (GET /order)

- Service ListOrders com GRPC

- Query ListOrders GraphQL

Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

# Important Note to use evans to call gRPC Server
- in the video, Wesley doen't show how to use evans to call gRPC Server, so I had to search for it. You have to use the following commands:
```
$ evans -r repl
$ shown package
$ package pb
$ show service
```