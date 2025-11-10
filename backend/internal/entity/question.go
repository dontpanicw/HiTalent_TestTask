package entity

import "time"

// Question - вопрос
type Question struct {
	Id        int       `gorm:"primaryKey;column:id" json:"id"`
	Text      string    `gorm:"column:text;not null" json:"text"` //(текст вопроса)
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Answers   []Answer  `gorm:"foreignKey:QuestionId;constraint:OnDelete:CASCADE" json:"answers,omitempty"`
}

func (Question) TableName() string {
	return "questions"
}
