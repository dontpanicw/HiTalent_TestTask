package postgres

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ repo.AnswerRepo = (*AnswerRepo)(nil)

type AnswerRepo struct {
	db *gorm.DB
}

func NewAnswerRepo(db *gorm.DB) *AnswerRepo {
	return &AnswerRepo{
		db: db,
	}
}

func (a *AnswerRepo) CreateAnswer(ctx context.Context, answer *entity.Answer) error {
	// Проверяем, существует ли вопрос
	var question entity.Question
	if err := a.db.WithContext(ctx).First(&question, answer.QuestionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}

	if err := a.db.WithContext(ctx).Create(answer).Error; err != nil {
		return err
	}
	return nil
}

func (a *AnswerRepo) GetAnswer(ctx context.Context, answerId int) (*entity.Answer, error) {
	var answer entity.Answer
	if err := a.db.WithContext(ctx).First(&answer, answerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("answer not found")
		}
		return nil, err
	}
	return &answer, nil
}

func (a *AnswerRepo) DeleteAnswer(ctx context.Context, answerId int) error {
	result := a.db.WithContext(ctx).Delete(&entity.Answer{}, answerId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("answer not found")
	}
	return nil
}
