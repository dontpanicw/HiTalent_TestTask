package server

import (
	"HiTalent_TestTask/backend/internal/cases"
	"HiTalent_TestTask/backend/internal/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Handlers struct {
	questionCase *cases.QuestionCase
	answerCase   *cases.AnswerCase
	logger       *zap.Logger
}

func NewHandlers(questionCase *cases.QuestionCase, answerCase *cases.AnswerCase, logger *zap.Logger) *Handlers {
	return &Handlers{
		questionCase: questionCase,
		answerCase:   answerCase,
		logger:       logger,
	}
}

// Question Handlers

func (h *Handlers) GetQuestionList(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questionCase.GetQuestionList(r.Context())
	if err != nil {
		h.logger.Error("Failed to get question list", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question entity.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if question.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	if err := h.questionCase.CreateQuestion(r.Context(), &question); err != nil {
		h.logger.Error("Failed to create question", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		return
	}
}

func (h *Handlers) GetQuestion(w http.ResponseWriter, r *http.Request, questionId int) {
	question, err := h.questionCase.GetQuestion(r.Context(), questionId)
	if err != nil {
		if err.Error() == "question not found" {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get question", zap.Int("id", questionId), zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		return
	}
}

func (h *Handlers) DeleteQuestion(w http.ResponseWriter, r *http.Request, questionId int) {
	if err := h.questionCase.DeleteQuestion(r.Context(), questionId); err != nil {
		if err.Error() == "question not found" {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete question", zap.Int("id", questionId), zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Answer Handlers

func (h *Handlers) CreateAnswer(w http.ResponseWriter, r *http.Request, questionId int) {
	var answer entity.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if answer.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	if answer.UserId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	answer.QuestionId = questionId

	if err := h.answerCase.CreateAnswer(r.Context(), &answer); err != nil {
		if err.Error() == "question not found" {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to create answer", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		return
	}
}

func (h *Handlers) GetAnswer(w http.ResponseWriter, r *http.Request, answerId int) {
	answer, err := h.answerCase.GetAnswer(r.Context(), answerId)
	if err != nil {
		if err.Error() == "answer not found" {
			http.Error(w, "Answer not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get answer", zap.Int("id", answerId), zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		return
	}
}

func (h *Handlers) DeleteAnswer(w http.ResponseWriter, r *http.Request, answerId int) {
	if err := h.answerCase.DeleteAnswer(r.Context(), answerId); err != nil {
		if err.Error() == "answer not found" {
			http.Error(w, "Answer not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete answer", zap.Int("id", answerId), zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractID извлекает ID из пути
func extractID(path string, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// parseInt преобразует строку в int
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	return strconv.Atoi(s)
}
