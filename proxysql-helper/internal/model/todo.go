package model

import (
	"time"
)

type Todo struct {
	ID        int64
	Title     string
	CreatedAt time.Time
}
