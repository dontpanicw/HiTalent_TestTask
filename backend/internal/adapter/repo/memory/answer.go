package memory

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"
	"errors"
	"sync"
	"time"
)

var _ repo.AnswerRepo = (*AnswerRepo)(nil)

type AnswerRepo struct {
	mu           sync.RWMutex
	answers      map[int]*entity.Answer
	nextID       int
	questionRepo *QuestionRepo // Для проверки существования вопроса
}

func NewAnswerRepo(questionRepo *QuestionRepo) *AnswerRepo {
	repo := &AnswerRepo{
		answers:      make(map[int]*entity.Answer),
		nextID:       1,
		questionRepo: questionRepo,
	}
	// Устанавливаем обратную ссылку для загрузки ответов
	questionRepo.SetAnswerRepo(repo)
	return repo
}

func (a *AnswerRepo) CreateAnswer(ctx context.Context, answer *entity.Answer) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Проверяем существование вопроса
	a.questionRepo.mu.RLock()
	_, exists := a.questionRepo.questions[answer.QuestionId]
	a.questionRepo.mu.RUnlock()

	if !exists {
		return errors.New("question not found")
	}

	answer.ID = a.nextID
	a.nextID++
	if answer.CreatedAt.IsZero() {
		answer.CreatedAt = time.Now()
	}
	a.answers[answer.ID] = answer
	return nil
}

func (a *AnswerRepo) GetAnswer(ctx context.Context, answerId int) (*entity.Answer, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	answer, exists := a.answers[answerId]
	if !exists {
		return nil, errors.New("answer not found")
	}

	result := *answer
	return &result, nil
}

func (a *AnswerRepo) DeleteAnswer(ctx context.Context, answerId int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.answers[answerId]; !exists {
		return errors.New("answer not found")
	}

	delete(a.answers, answerId)
	return nil
}

// SetAnswerForTesting устанавливает ответ для тестирования
func (a *AnswerRepo) SetAnswerForTesting(answer *entity.Answer) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.answers[answer.ID] = answer
	if answer.ID >= a.nextID {
		a.nextID = answer.ID + 1
	}
}
