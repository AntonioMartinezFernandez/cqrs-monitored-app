package di

import (
	"context"
	"os"
	"os/signal"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/configs"

	pkg_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_command_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/command"
	pkg_event_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/event"
	pkg_query_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/query"
	pkg_json_schema "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-schema"
	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	pkg_mutex_service "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/mutex-service"
	pkg_observability "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/observability"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"

	"github.com/joho/godotenv"
)

type CommonServices struct {
	Config      configs.Config
	Environment configs.Environment

	Logger              pkg_logger.Logger
	DistributedMutex    pkg_mutex_service.MutexService
	JsonSchemaValidator *pkg_json_schema.JsonSchemaValidator
	Observability       *pkg_observability.OtelObservability
	UlidProvider        pkg_utils.UlidProvider
	UuidProvider        pkg_utils.UuidProvider
	TimeProvider        pkg_utils.DateTimeProvider
	CommandBus          *pkg_command_bus.CommandBus
	QueryBus            *pkg_query_bus.QueryBus
	EventBus            *pkg_event_bus.EventBus
}

func InitCommonServices(ctx context.Context) *CommonServices {
	config := initConfig()
	environment := configs.NewEnvironmentFromRawEnvVar(config.AppEnv)
	logger := pkg_logger.NewOtelInstrumentalizedLogger(config.LogLevel)
	distributedMutex := pkg_mutex_service.NewInmemoryMutexService(logger)
	jsonSchemaValidator := pkg_json_schema.NewJsonSchemaValidator(config.JsonSchemaBasePath)
	ulidProvider := pkg_utils.NewRandomUlidProvider()
	uuidProvider := pkg_utils.NewRandomUuidProvider()
	timeProvider := pkg_utils.NewSystemTimeProvider()
	commandBus := pkg_command_bus.InitCommandBus(logger, distributedMutex)
	queryBus := pkg_query_bus.InitQueryBus(logger)
	eventBus := pkg_event_bus.NewEventBus()

	grpcConnection, grpcErr := pkg_observability.InitGrpcConnInsecure(config.OtelGrpcHost, config.OtelGrpcPort)
	if grpcErr != nil {
		logger.Error(ctx, "error establishing grpc connection with OTEL collector")
	}
	otelObservability, obsErr := pkg_observability.InitOpenTelemetryObservability(
		ctx,
		grpcConnection,
		config.AppServiceName,
		config.AppVersion,
	)
	if obsErr != nil {
		logger.Error(ctx, "error initializing OpenTelemetry observability")
	}

	return &CommonServices{
		Config:      config,
		Environment: environment,

		Logger:              logger,
		DistributedMutex:    distributedMutex,
		JsonSchemaValidator: &jsonSchemaValidator,
		Observability:       otelObservability,
		UlidProvider:        ulidProvider,
		UuidProvider:        uuidProvider,
		TimeProvider:        timeProvider,
		CommandBus:          commandBus,
		QueryBus:            queryBus,
		EventBus:            eventBus,
	}
}

/* HELPERS */

func InitCommonServicesWithEnvFiles(envFiles ...string) *CommonServices {
	ctx := context.Background()
	err := godotenv.Overload(envFiles...)
	if err != nil {
		panic(err)
	}

	return InitCommonServices(ctx)
}

func RootContext() (context.Context, context.CancelFunc) {
	rootCtx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, os.Kill,
	)
	return rootCtx, cancel
}

func initConfig() configs.Config {
	return configs.LoadEnvConfig()
}

func registerQueryOrPanic(
	queryBus pkg_query_bus.Bus,
	query pkg_bus.Dto,
	handler pkg_query_bus.QueryHandler,
) {
	if err := queryBus.RegisterQuery(query, handler); err != nil {
		panic(err)
	}
}

func registerCommandOrPanic(
	commandBus pkg_command_bus.Bus,
	cmd pkg_bus.Dto,
	handler pkg_command_bus.CommandHandler,
) {
	if err := commandBus.RegisterCommand(cmd, handler); err != nil {
		panic(err)
	}
}

func registerEvent(
	eventBus pkg_event_bus.EventBus,
	eventName string,
	eventHandler pkg_event_bus.EventHandler,
) {
	eventBus.Subscribe(eventName, eventHandler)
}
