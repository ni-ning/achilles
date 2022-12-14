package model

import "time"

type Todo struct {
	ID        uint64    `gorm:"primarykey;" json:"id"`
	Title     string    `gorm:"index" json:"title"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `json:"-"`
}
