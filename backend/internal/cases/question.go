package cases

import (
	"HiTalent_TestTask/backend/internal/entity"
	"HiTalent_TestTask/backend/internal/port/repo"
	"context"

	"go.uber.org/zap"
)

type QuestionCase struct {
	questionRepo repo.QuestionRepo
	logger       *zap.Logger
}

func NewQuestionCase(questionRepo repo.QuestionRepo, logger *zap.Logger) *QuestionCase {
	return &QuestionCase{
		questionRepo: questionRepo,
		logger:       logger,
	}
}

func (q *QuestionCase) GetQuestionList(ctx context.Context) (*[]entity.Question, error) {
	q.logger.Info("Getting question list")
	questions, err := q.questionRepo.GetQuestionList(ctx)
	if err != nil {
		q.logger.Error("Failed to get question list", zap.Error(err))
		return nil, err
	}
	return questions, nil
}

func (q *QuestionCase) CreateQuestion(ctx context.Context, question *entity.Question) error {
	q.logger.Info("Creating question", zap.String("text", question.Text))
	if err := q.questionRepo.CreateQuestion(ctx, question); err != nil {
		q.logger.Error("Failed to create question", zap.Error(err))
		return err
	}
	q.logger.Info("Question created successfully", zap.Int("id", question.Id))
	return nil
}

func (q *QuestionCase) GetQuestion(ctx context.Context, questionId int) (*entity.Question, error) {
	q.logger.Info("Getting question", zap.Int("id", questionId))
	question, err := q.questionRepo.GetQuestion(ctx, questionId)
	if err != nil {
		q.logger.Error("Failed to get question", zap.Int("id", questionId), zap.Error(err))
		return nil, err
	}
	return question, nil
}

func (q *QuestionCase) DeleteQuestion(ctx context.Context, questionId int) error {
	q.logger.Info("Deleting question", zap.Int("id", questionId))
	if err := q.questionRepo.DeleteQuestion(ctx, questionId); err != nil {
		q.logger.Error("Failed to delete question", zap.Int("id", questionId), zap.Error(err))
		return err
	}
	q.logger.Info("Question deleted successfully", zap.Int("id", questionId))
	return nil
}
