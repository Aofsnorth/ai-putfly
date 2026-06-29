package main

import (
	"log/slog"

	"github.com/user/aofsnorth/ai-putfly/internal/config"
	"github.com/user/aofsnorth/ai-putfly/internal/handler"
)

func main() {
	conf, err := config.ReadConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()
	srv.Listen()
	for p := range srv.Accept() {
		p.Handle(handler.NewListenChat()) // internal/handler
	}
}
