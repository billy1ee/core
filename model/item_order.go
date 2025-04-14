package model

import (
	"time"
)

// Item represents a product or service that can be ordered
type Item struct {
	ID          string    `json:"id" bson:"_id"`
	AccountID   string    `json:"accountId" bson:"accountId"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Created     time.Time `json:"created" bson:"created"`
}

// Order represents a customer order containing multiple items
type Order struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"accountId" bson:"accountId"`
	UserID    string    `json:"userId" bson:"userId"`
	Items     []string  `json:"items" bson:"items"`
	Total     float64   `json:"total" bson:"total"`
	Status    string    `json:"status" bson:"status"`
	Created   time.Time `json:"created" bson:"created"`
}