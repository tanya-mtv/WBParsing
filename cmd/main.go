package main

import (
	"flag"
	"log"
	"parsingWB/internal/agent"
	"parsingWB/internal/config"
	"parsingWB/internal/logger"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	ag := agent.NewAgent(cfg, appLogger)

	if err := ag.Run(); err != nil {
		appLogger.Fatal(err)
	}

}
