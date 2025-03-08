package configs

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	AppServiceName string `env:"APP_SERVICE_NAME"`
	AppEnv         string `env:"APP_ENV"`
	AppVersion     string `env:"APP_VERSION"`

	LogLevel           string `env:"LOG_LEVEL"`
	JsonSchemaBasePath string `env:"JSON_SCHEMA_BASE_PATH"`

	HttpHost         string `env:"HTTP_HOST"`
	HttpPort         string `env:"HTTP_PORT"`
	HttpReadTimeout  int    `env:"HTTP_READ_TIMEOUT"`
	HttpWriteTimeout int    `env:"HTTP_WRITE_TIMEOUT"`

	OtelGrpcHost string `env:"OTEL_GRPC_HOST"`
	OtelGrpcPort string `env:"OTEL_GRPC_PORT"`
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
