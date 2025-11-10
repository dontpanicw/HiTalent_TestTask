package app

import (
	"HiTalent_TestTask/backend/config"
	"HiTalent_TestTask/backend/internal/adapter/repo/postgres"
	"HiTalent_TestTask/backend/internal/cases"
	"HiTalent_TestTask/backend/internal/input/http/server"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Start(cfg config.Config, logger *zap.Logger) error {
	// Создаем подключение к БД через GORM
	db, err := postgres.NewGormDB(cfg.PgConnStr)
	if err != nil {
		logger.Fatal("Failed to connect to postgres", zap.Error(err))
		return fmt.Errorf("failed to create gorm db: %w", err)
	}

	// Создаем репозитории
	questionRepo := postgres.NewQuestionRepo(db)
	answerRepo := postgres.NewAnswerRepo(db)

	// Создаем cases (бизнес-логика)
	questionCase := cases.NewQuestionCase(questionRepo, logger)
	answerCase := cases.NewAnswerCase(answerRepo, logger)

	// Создаем HTTP сервер
	srv := server.NewServer(questionCase, answerCase, logger)

	logger.Info("Starting server", zap.String("port", cfg.HTTPPort))
	return http.ListenAndServe(cfg.HTTPPort, srv)
}
