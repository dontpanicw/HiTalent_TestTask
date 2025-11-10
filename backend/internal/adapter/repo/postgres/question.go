package postgres

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ repo.QuestionRepo = (*QuestionRepo)(nil)

type QuestionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepo {
	return &QuestionRepo{
		db: db,
	}
}

func (q *QuestionRepo) GetQuestionList(ctx context.Context) (*[]entity.Question, error) {
	var questions []entity.Question
	if err := q.db.WithContext(ctx).Find(&questions).Error; err != nil {
		return nil, err
	}
	return &questions, nil
}

func (q *QuestionRepo) CreateQuestion(ctx context.Context, question *entity.Question) error {
	if err := q.db.WithContext(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

func (q *QuestionRepo) GetQuestion(ctx context.Context, questionId int) (*entity.Question, error) {
	var question entity.Question
	if err := q.db.WithContext(ctx).Preload("Answers").First(&question, questionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

func (q *QuestionRepo) DeleteQuestion(ctx context.Context, questionId int) error {
	result := q.db.WithContext(ctx).Delete(&entity.Question{}, questionId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("question not found")
	}
	return nil
}
