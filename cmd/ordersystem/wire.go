//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/entity"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/event"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/database"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/infra/web"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/internal/usecase"
	"github.com/cesarmiggiolaro/go/cleanarchitecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}

func NewGetOrdersUseCase(db *sql.DB) *usecase.GetOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewGetOrdersUseCase,
	)
	return &usecase.GetOrdersUseCase{}
}

func GetWebOrderHandler(db *sql.DB) *web.WebGetOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		web.GetwWebOrderHandler,
	)
	return &web.WebGetOrderHandler{}
}
