package cases

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"

	"go.uber.org/zap"
)

type AnswerCase struct {
	answerRepo repo.AnswerRepo
	logger     *zap.Logger
}

func NewAnswerCase(answerRepo repo.AnswerRepo, logger *zap.Logger) *AnswerCase {
	return &AnswerCase{
		answerRepo: answerRepo,
		logger:     logger,
	}
}

func (a *AnswerCase) CreateAnswer(ctx context.Context, answer *entity.Answer) error {
	a.logger.Info("Creating answer",
		zap.Int("question_id", answer.QuestionId),
		zap.String("user_id", answer.UserId))

	if err := a.answerRepo.CreateAnswer(ctx, answer); err != nil {
		a.logger.Error("Failed to create answer", zap.Error(err))
		return err
	}

	a.logger.Info("Answer created successfully", zap.Int("id", answer.ID))
	return nil
}

func (a *AnswerCase) GetAnswer(ctx context.Context, answerId int) (*entity.Answer, error) {
	a.logger.Info("Getting answer", zap.Int("id", answerId))
	answer, err := a.answerRepo.GetAnswer(ctx, answerId)
	if err != nil {
		a.logger.Error("Failed to get answer", zap.Int("id", answerId), zap.Error(err))
		return nil, err
	}
	return answer, nil
}

func (a *AnswerCase) DeleteAnswer(ctx context.Context, answerId int) error {
	a.logger.Info("Deleting answer", zap.Int("id", answerId))
	if err := a.answerRepo.DeleteAnswer(ctx, answerId); err != nil {
		a.logger.Error("Failed to delete answer", zap.Int("id", answerId), zap.Error(err))
		return err
	}
	a.logger.Info("Answer deleted successfully", zap.Int("id", answerId))
	return nil
}
