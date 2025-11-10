package server

import (
	"HiTalent_TestTask/backend/internal/cases"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	mux    *http.ServeMux
	logger *zap.Logger
}

func NewServer(questionCase *cases.QuestionCase, answerCase *cases.AnswerCase, logger *zap.Logger) *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		logger: logger,
	}

	handlers := NewHandlers(questionCase, answerCase, logger)

	// Регистрируем обработчики
	s.mux.HandleFunc("/questions/", s.questionsHandler(handlers))
	s.mux.HandleFunc("/answers/", s.answersHandler(handlers))

	return s
}

// questionsHandler обрабатывает все запросы к /questions/
func (s *Server) questionsHandler(h *Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path

		// Проверяем специальный случай: POST /questions/{id}/answers/
		if strings.Contains(fullPath, "/answers") && r.Method == http.MethodPost {
			questionID, err := s.extractQuestionIDFromAnswerPath(fullPath)
			if err != nil {
				http.Error(w, "Invalid question ID", http.StatusBadRequest)
				return
			}
			h.CreateAnswer(w, r, questionID)
			return
		}

		// Обрабатываем остальные пути
		path := strings.TrimPrefix(fullPath, "/questions")
		path = strings.Trim(path, "/")

		if path == "" {
			// Путь /questions/ или /questions
			switch r.Method {
			case http.MethodGet:
				// GET /questions/
				h.GetQuestionList(w, r)
			case http.MethodPost:
				// POST /questions/
				h.CreateQuestion(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Путь содержит ID: /questions/{id}
		questionID, err := s.extractID(path)
		if err != nil {
			http.Error(w, "Invalid question ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// GET /questions/{id}
			h.GetQuestion(w, r, questionID)
		case http.MethodDelete:
			// DELETE /questions/{id}
			h.DeleteQuestion(w, r, questionID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// answersHandler обрабатывает все запросы к /answers/
func (s *Server) answersHandler(h *Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/answers")
		path = strings.Trim(path, "/")

		answerID, err := s.extractID(path)
		if err != nil {
			http.Error(w, "Invalid answer ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// GET /answers/{id}
			h.GetAnswer(w, r, answerID)

		case http.MethodDelete:
			// DELETE /answers/{id}
			h.DeleteAnswer(w, r, answerID)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// extractID извлекает ID из пути
func (s *Server) extractID(path string) (int, error) {
	if path == "" {
		return 0, http.ErrMissingFile
	}
	// Берем первую часть пути как ID
	parts := strings.Split(path, "/")
	return parseInt(parts[0])
}

// extractQuestionIDFromAnswerPath извлекает ID вопроса из пути вида /questions/{id}/answers/ или /questions/{id}/answers
func (s *Server) extractQuestionIDFromAnswerPath(fullPath string) (int, error) {
	// Убираем префикс /questions/
	path := strings.TrimPrefix(fullPath, "/questions/")
	// Убираем суффикс /answers/ или /answers
	path = strings.TrimSuffix(path, "/answers/")
	path = strings.TrimSuffix(path, "/answers")
	path = strings.Trim(path, "/")

	// Берем первую часть как ID вопроса
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, http.ErrMissingFile
	}

	return parseInt(parts[0])
}

// ServeHTTP реализует http.Handler с middleware для логирования и recovery
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Recovery middleware
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Panic recovered",
				zap.Any("error", err),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Обертка для ResponseWriter для отслеживания статус-кода
	wrapped := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Обрабатываем запрос
	s.mux.ServeHTTP(wrapped, r)

	// Логирование после обработки
	duration := time.Since(start)
	s.logger.Info("HTTP request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Int("status", wrapped.statusCode),
		zap.Duration("duration", duration),
	)
}

// responseWriter обертка для ResponseWriter для отслеживания статус-кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
