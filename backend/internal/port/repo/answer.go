package repo

import (
	"HiTalent_TestTask/backend/internal/entity"
	"context"
)

type AnswerRepo interface {
	CreateAnswer(ctx context.Context, answer *entity.Answer) error
	GetAnswer(ctx context.Context, answerId int) (*entity.Answer, error)
	DeleteAnswer(ctx context.Context, answerId int) error
}

//POST /questions/{id}/answers/ — добавить ответ к вопросу
//GET /answers/{id} — получить конкретный ответ
//DELETE /answers/{id} — удалить ответ
