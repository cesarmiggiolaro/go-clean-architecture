package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/configs"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/event/handler"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/graph"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/grpc/pb"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/grpc/service"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/web/webserver"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	//REST
	webserver := webserver.NewWebServer(configs.WebServerPort)

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)

	getOrdersUseCase := NewGetOrdersUseCase(db)
	webGetOrderHandler := GetWebOrderHandler(db)
	webserver.AddHandler("/orders", webGetOrderHandler.Get)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	//gRPC
	orderService := service.NewOrderService(*createOrderUseCase, *getOrdersUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	go grpcServer.Serve(lis)

	//fmt.Printf("%v", grpcServer.GetServiceInfo())

	//GraphQL
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		GetOrdersUseCase:   *getOrdersUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
