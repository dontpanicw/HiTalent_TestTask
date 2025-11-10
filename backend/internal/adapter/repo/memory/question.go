package memory

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"
	"errors"
	"sync"
	"time"
)

var _ repo.QuestionRepo = (*QuestionRepo)(nil)

type QuestionRepo struct {
	mu         sync.RWMutex
	questions  map[int]*entity.Question
	nextID     int
	answerRepo *AnswerRepo // Для загрузки ответов
}

func (q *QuestionRepo) SetAnswerRepo(answerRepo *AnswerRepo) {
	q.answerRepo = answerRepo
}

func NewQuestionRepo() *QuestionRepo {
	return &QuestionRepo{
		questions: make(map[int]*entity.Question),
		nextID:    1,
	}
}

func (q *QuestionRepo) GetQuestionList(ctx context.Context) (*[]entity.Question, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	questions := make([]entity.Question, 0, len(q.questions))
	for _, question := range q.questions {
		questions = append(questions, *question)
	}
	return &questions, nil
}

func (q *QuestionRepo) CreateQuestion(ctx context.Context, question *entity.Question) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	question.Id = q.nextID
	q.nextID++
	if question.CreatedAt.IsZero() {
		question.CreatedAt = time.Now()
	}
	q.questions[question.Id] = question
	return nil
}

func (q *QuestionRepo) GetQuestion(ctx context.Context, questionId int) (*entity.Question, error) {
	q.mu.RLock()
	question, exists := q.questions[questionId]
	if !exists {
		q.mu.RUnlock()
		return nil, errors.New("question not found")
	}

	// Копируем вопрос
	result := *question
	result.Answers = []entity.Answer{}
	q.mu.RUnlock()

	// Загружаем ответы, если answerRepo установлен
	if q.answerRepo != nil {
		q.answerRepo.mu.RLock()
		answers := make([]entity.Answer, 0)
		for _, answer := range q.answerRepo.answers {
			if answer.QuestionId == questionId {
				answers = append(answers, *answer)
			}
		}
		q.answerRepo.mu.RUnlock()
		result.Answers = answers
	}

	return &result, nil
}

func (q *QuestionRepo) DeleteQuestion(ctx context.Context, questionId int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.questions[questionId]; !exists {
		return errors.New("question not found")
	}

	delete(q.questions, questionId)
	return nil
}

// SetQuestionForTesting устанавливает вопрос для тестирования
func (q *QuestionRepo) SetQuestionForTesting(question *entity.Question) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.questions[question.Id] = question
	if question.Id >= q.nextID {
		q.nextID = question.Id + 1
	}
}
