package models

type Account struct {
	BaseModel

	DocumentNumber string `gorm:"not null;uniqueIndex"`
}
