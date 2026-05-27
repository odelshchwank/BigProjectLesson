package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/odelshchwank/BigProjectLesson/internal/core/config"
	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_pgx_pool "github.com/odelshchwank/BigProjectLesson/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/middleware"
	core_http_server "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/odelshchwank/BigProjectLesson/internal/features/tasks/repository/postgres"
	tasks_service "github.com/odelshchwank/BigProjectLesson/internal/features/tasks/service"
	tasks_transport_http "github.com/odelshchwank/BigProjectLesson/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/odelshchwank/BigProjectLesson/internal/features/users/repository/postgres"
	users_service "github.com/odelshchwank/BigProjectLesson/internal/features/users/service"
	users_transport_http "github.com/odelshchwank/BigProjectLesson/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	// Инициализация логгера
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()
	// Закончили инициализацию логгера

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	// Инициализация бд-пула
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()
	// Закончили инициализацию бд-пула

	// Инициализируем все наши слои (бд, сервис и транспорт) фичи users
	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)
	// Закончили инициализацию слоёв фичи users

	// Инициализируем все наши слои (бд, сервис и транспорт) фичи tasks
	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)
	// Закончили инициализацию слоёв фичи tasks

	// Инициализируем сервер
	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)

	/* apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion2,
		core_http_middleware.Dummy("api v2 middleware"),
	)
	apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...) */

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
