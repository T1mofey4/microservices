package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/T1mofey4/microservices/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load env: %w", err)
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	paymentGRPCcfg, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerConfig,
		PaymentGRPC: paymentGRPCcfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
