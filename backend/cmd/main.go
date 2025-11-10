package main

import (
	"HiTalent_TestTask/backend/config"
	"HiTalent_TestTask/backend/internal/app"

	"go.uber.org/zap"
)

func main() {
	//create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logger.Sync()

	cfg, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatal("error creating config", zap.Error(err))
	}

	// Миграции применяются автоматически в app.Start через NewGormDB
	if err := app.Start(cfg, logger); err != nil {
		logger.Fatal("failed to start application", zap.Error(err))
	}
}
