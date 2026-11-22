package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentGRPCConfig struct {
	raw paymentGRPCEnvConfig
}

func NewPaymentGrpcConfig() (*paymentGRPCConfig, error) {
	var raw paymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("grpc parse config failed: %w", err)
	}

	return &paymentGRPCConfig{raw: raw}, nil
}

func (cfg *paymentGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
