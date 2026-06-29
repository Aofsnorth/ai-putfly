package server

import (
	"log/slog"

	"github.com/df-mc/dragonfly/server"
	"github.com/user/aofsnorth/ai-putfly/internal/config"
)

func Setup() *server.Server {
	conf, err := config.ReadConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	return srv
}
