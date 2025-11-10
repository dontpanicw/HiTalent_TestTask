package repo

import (
	"HiTalent_TestTask/backend/internal/entity"
	"context"
)

type QuestionRepo interface {
	GetQuestionList(ctx context.Context) (*[]entity.Question, error)
	CreateQuestion(ctx context.Context, question *entity.Question) error
	GetQuestion(ctx context.Context, questionId int) (*entity.Question, error)
	DeleteQuestion(ctx context.Context, questionId int) error
}

//GET /questions/ — список всех вопросов
//POST /questions/ — создать новый вопрос
//GET /questions/{id} — получить вопрос и все ответы на него
//DELETE /questions/{id} — удалить вопрос (вместе с ответами)
