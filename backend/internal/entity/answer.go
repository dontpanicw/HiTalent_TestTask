package entity

import "time"

type Answer struct {
	ID         int       `gorm:"primaryKey;column:id" json:"id"`
	QuestionId int       `gorm:"column:question_id;not null;index" json:"question_id"`
	UserId     string    `gorm:"column:user_id;not null;index" json:"user_id"` //uuid
	Text       string    `gorm:"column:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Question   Question  `gorm:"foreignKey:QuestionId" json:"question,omitempty"`
}

func (Answer) TableName() string {
	return "answers"
}
