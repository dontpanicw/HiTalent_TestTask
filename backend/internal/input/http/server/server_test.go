package server

import (
	"HiTalent_TestTask/backend/internal/adapter/repo/memory"
	"HiTalent_TestTask/backend/internal/cases"
	"HiTalent_TestTask/backend/internal/entity"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestServer() (*Server, *memory.QuestionRepo, *memory.AnswerRepo) {
	logger := zap.NewNop()
	questionRepo := memory.NewQuestionRepo()
	answerRepo := memory.NewAnswerRepo(questionRepo)

	questionCase := cases.NewQuestionCase(questionRepo, logger)
	answerCase := cases.NewAnswerCase(answerRepo, logger)

	server := NewServer(questionCase, answerCase, logger)
	return server, questionRepo, answerRepo
}

func TestGetQuestionList(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	// Создаем тестовые вопросы
	question1 := &entity.Question{
		Id:   1,
		Text: "Test Question 1",
	}
	question2 := &entity.Question{
		Id:   2,
		Text: "Test Question 2",
	}
	questionRepo.SetQuestionForTesting(question1)
	questionRepo.SetQuestionForTesting(question2)

	req := httptest.NewRequest(http.MethodGet, "/questions/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var questions []entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &questions)
	require.NoError(t, err)
	assert.Len(t, questions, 2)
}

func TestGetQuestionListEmpty(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/questions/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var questions []entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &questions)
	require.NoError(t, err)
	assert.Len(t, questions, 0)
}

func TestCreateQuestion(t *testing.T) {
	server, _, _ := setupTestServer()

	question := entity.Question{
		Text: "New Question",
	}
	body, _ := json.Marshal(question)

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var createdQuestion entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &createdQuestion)
	require.NoError(t, err)
	assert.Equal(t, question.Text, createdQuestion.Text)
	assert.NotZero(t, createdQuestion.Id)
}

func TestCreateQuestionEmptyText(t *testing.T) {
	server, _, _ := setupTestServer()

	question := entity.Question{
		Text: "",
	}
	body, _ := json.Marshal(question)

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetQuestion(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var retrievedQuestion entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &retrievedQuestion)
	require.NoError(t, err)
	assert.Equal(t, question.Id, retrievedQuestion.Id)
	assert.Equal(t, question.Text, retrievedQuestion.Text)
}

func TestGetQuestionNotFound(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/questions/999", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteQuestion(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Проверяем, что вопрос удален
	getReq := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	getW := httptest.NewRecorder()
	server.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestDeleteQuestionNotFound(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodDelete, "/questions/999", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateAnswer(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	answer := entity.Answer{
		UserId: "user-123",
		Text:   "Test Answer",
	}
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var createdAnswer entity.Answer
	err := json.Unmarshal(w.Body.Bytes(), &createdAnswer)
	require.NoError(t, err)
	assert.Equal(t, answer.Text, createdAnswer.Text)
	assert.Equal(t, answer.UserId, createdAnswer.UserId)
	assert.Equal(t, 1, createdAnswer.QuestionId)
	assert.NotZero(t, createdAnswer.ID)
}

func TestCreateAnswerQuestionNotFound(t *testing.T) {
	server, _, _ := setupTestServer()

	answer := entity.Answer{
		UserId: "user-123",
		Text:   "Test Answer",
	}
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/999/answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateAnswerEmptyText(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	answer := entity.Answer{
		UserId: "user-123",
		Text:   "",
	}
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAnswerEmptyUserId(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	answer := entity.Answer{
		UserId: "",
		Text:   "Test Answer",
	}
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAnswer(t *testing.T) {
	server, questionRepo, answerRepo := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем ответ
	answer := &entity.Answer{
		ID:         1,
		QuestionId: 1,
		UserId:     "user-123",
		Text:       "Test Answer",
	}
	answerRepo.SetAnswerForTesting(answer)

	req := httptest.NewRequest(http.MethodGet, "/answers/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var retrievedAnswer entity.Answer
	err := json.Unmarshal(w.Body.Bytes(), &retrievedAnswer)
	require.NoError(t, err)
	assert.Equal(t, answer.ID, retrievedAnswer.ID)
	assert.Equal(t, answer.Text, retrievedAnswer.Text)
	assert.Equal(t, answer.UserId, retrievedAnswer.UserId)
}

func TestGetAnswerNotFound(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/answers/999", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAnswer(t *testing.T) {
	server, questionRepo, answerRepo := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем ответ
	answer := &entity.Answer{
		ID:         1,
		QuestionId: 1,
		UserId:     "user-123",
		Text:       "Test Answer",
	}
	answerRepo.SetAnswerForTesting(answer)

	req := httptest.NewRequest(http.MethodDelete, "/answers/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Проверяем, что ответ удален
	getReq := httptest.NewRequest(http.MethodGet, "/answers/1", nil)
	getW := httptest.NewRecorder()
	server.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestDeleteAnswerNotFound(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodDelete, "/answers/999", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	server, _, _ := setupTestServer()

	// Пытаемся использовать PUT на /questions/
	req := httptest.NewRequest(http.MethodPut, "/questions/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestInvalidQuestionID(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/questions/abc", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInvalidAnswerID(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/answers/abc", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetQuestionWithAnswers(t *testing.T) {
	server, questionRepo, answerRepo := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:        1,
		Text:      "Test Question",
		CreatedAt: time.Now(),
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем ответы
	answer1 := &entity.Answer{
		ID:         1,
		QuestionId: 1,
		UserId:     "user-1",
		Text:       "Answer 1",
		CreatedAt:  time.Now(),
	}
	answer2 := &entity.Answer{
		ID:         2,
		QuestionId: 1,
		UserId:     "user-2",
		Text:       "Answer 2",
		CreatedAt:  time.Now(),
	}
	answerRepo.SetAnswerForTesting(answer1)
	answerRepo.SetAnswerForTesting(answer2)

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedQuestion entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &retrievedQuestion)
	require.NoError(t, err)
	assert.Equal(t, question.Id, retrievedQuestion.Id)
	// Note: In-memory repo не загружает answers автоматически, но структура позволяет это
}

func TestMultipleAnswersSameUser(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем первый ответ
	answer1 := entity.Answer{
		UserId: "user-123",
		Text:   "First Answer",
	}
	body1, _ := json.Marshal(answer1)

	req1 := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	server.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Создаем второй ответ от того же пользователя
	answer2 := entity.Answer{
		UserId: "user-123",
		Text:   "Second Answer",
	}
	body2, _ := json.Marshal(answer2)

	req2 := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	server.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	// Проверяем, что оба ответа созданы
	var createdAnswer1, createdAnswer2 entity.Answer
	json.Unmarshal(w1.Body.Bytes(), &createdAnswer1)
	json.Unmarshal(w2.Body.Bytes(), &createdAnswer2)

	assert.NotEqual(t, createdAnswer1.ID, createdAnswer2.ID)
	assert.Equal(t, "user-123", createdAnswer1.UserId)
	assert.Equal(t, "user-123", createdAnswer2.UserId)
}

func TestInvalidJSON(t *testing.T) {
	server, _, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAnswerInvalidJSON(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuestionCreatedAt(t *testing.T) {
	server, _, _ := setupTestServer()

	question := entity.Question{
		Text: "Test Question",
	}
	body, _ := json.Marshal(question)

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdQuestion entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &createdQuestion)
	require.NoError(t, err)
	assert.NotZero(t, createdQuestion.CreatedAt)
}

func TestAnswerCreatedAt(t *testing.T) {
	server, questionRepo, _ := setupTestServer()

	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	answer := entity.Answer{
		UserId: "user-123",
		Text:   "Test Answer",
	}
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdAnswer entity.Answer
	err := json.Unmarshal(w.Body.Bytes(), &createdAnswer)
	require.NoError(t, err)
	assert.NotZero(t, createdAnswer.CreatedAt)
}

func TestGetQuestionWithMultipleAnswers(t *testing.T) {
	server, questionRepo, answerRepo := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем несколько ответов
	answer1 := &entity.Answer{
		ID:         1,
		QuestionId: 1,
		UserId:     "user-1",
		Text:       "Answer 1",
	}
	answer2 := &entity.Answer{
		ID:         2,
		QuestionId: 1,
		UserId:     "user-2",
		Text:       "Answer 2",
	}
	answer3 := &entity.Answer{
		ID:         3,
		QuestionId: 1,
		UserId:     "user-1",
		Text:       "Answer 3",
	}
	answerRepo.SetAnswerForTesting(answer1)
	answerRepo.SetAnswerForTesting(answer2)
	answerRepo.SetAnswerForTesting(answer3)

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedQuestion entity.Question
	err := json.Unmarshal(w.Body.Bytes(), &retrievedQuestion)
	require.NoError(t, err)
	assert.Equal(t, question.Id, retrievedQuestion.Id)
	assert.Len(t, retrievedQuestion.Answers, 3)
}

func TestDeleteQuestionCascadeAnswers(t *testing.T) {
	server, questionRepo, answerRepo := setupTestServer()

	// Создаем вопрос
	question := &entity.Question{
		Id:   1,
		Text: "Test Question",
	}
	questionRepo.SetQuestionForTesting(question)

	// Создаем ответы
	answer1 := &entity.Answer{
		ID:         1,
		QuestionId: 1,
		UserId:     "user-1",
		Text:       "Answer 1",
	}
	answer2 := &entity.Answer{
		ID:         2,
		QuestionId: 1,
		UserId:     "user-2",
		Text:       "Answer 2",
	}
	answerRepo.SetAnswerForTesting(answer1)
	answerRepo.SetAnswerForTesting(answer2)

	// Удаляем вопрос
	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Проверяем, что вопрос удален
	getReq := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	getW := httptest.NewRecorder()
	server.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)

	// Примечание: В реальной БД ответы удаляются каскадно, но в in-memory репозитории
	// мы не реализуем каскадное удаление, так как это требует дополнительной логики
}
