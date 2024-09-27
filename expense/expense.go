package expense

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Expense struct {
	Id     int            `json:"id"`
	Title  string         `json:"title"`
	Amount float64        `json:"amount"`
	Note   string         `json:"note"`
	Tags   pq.StringArray `gorm:"type:text[]" json:"tags"`
}

type handler struct {
	DB *gorm.DB
}

func ExpenseHandler(db *gorm.DB) *handler {
	return &handler{db}
}

type Err struct {
	Message string `json:"message"`
}
