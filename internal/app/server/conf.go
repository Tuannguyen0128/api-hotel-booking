package server

import (
	"time"

	"github.com/gin-contrib/cors"
)

type (
	Config struct {
		Port                 uint          `mapstructure:"Port"`
		GracefulShutdownTime time.Duration `mapstructure:"GracefulShutdownTime"`
		ReadTimeout          time.Duration `mapstructure:"ReadTimeout"`
		WriteTimeout         time.Duration `mapstructure:"WriteTimeout"`
		IdleTimeout          time.Duration `mapstructure:"IdleTimeout"`
		Https                TLSConfig     `mapstructure:"HTTPS"`
		CORS                 cors.Config   `mapstructure:"CORS"`
	}

	TLSConfig struct {
		Enabled  bool   `mapstructure:"Enabled"`
		CertFile string `mapstructure:"CertFile"`
		KeyFile  string `mapstructure:"KeyFile"`
	}
)
