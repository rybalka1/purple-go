package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
	Price       float64        `json:"price"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
