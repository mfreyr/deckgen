package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mfreyr/deckgen/internal/config"
	"github.com/mfreyr/deckgen/internal/handler"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to the config file")
	dumpConfig := flag.Bool("dump-config", false, "dump the default config")
	flag.Parse()

	if *dumpConfig {
		if err := config.Dump(); err != nil {
			log.Fatalf("dump config error: %s\n", err)
		}
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config error load: %s\n", err)
	}

	handler := handler.New()

	server := &http.Server{
		Addr:           fmt.Sprintf("127.0.0.1:%d", cfg.Server.Port),
		Handler:        handler,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	run(cfg.Logger, server)
}
