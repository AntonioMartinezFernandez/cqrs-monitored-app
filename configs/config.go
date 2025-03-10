package configs

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	AppServiceName string `env:"APP_SERVICE_NAME,default=cqrs-monitored-app"`
	AppEnv         string `env:"APP_ENV,default=development"`
	AppVersion     string `env:"APP_VERSION,default=1.0.0"`

	LogLevel           string `env:"LOG_LEVEL,default=debug"`
	JsonSchemaBasePath string `env:"JSON_SCHEMA_BASE_PATH,default=./schema"`

	HttpHost         string `env:"HTTP_HOST,default=0.0.0.0"`
	HttpPort         string `env:"HTTP_PORT,default=8080"`
	HttpReadTimeout  int    `env:"HTTP_READ_TIMEOUT,default=30"`
	HttpWriteTimeout int    `env:"HTTP_WRITE_TIMEOUT,default=30"`

	OtelGrpcHost string `env:"OTEL_GRPC_HOST,default=otel-collector-opentelemetry-collector.otel-collector.svc.cluster.local"`
	OtelGrpcPort string `env:"OTEL_GRPC_PORT,default=4317"`
}

func LoadEnvConfig() Config {
	var config Config

	_ = godotenv.Load(".env")
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
