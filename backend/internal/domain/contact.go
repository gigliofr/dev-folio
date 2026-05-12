package domain

import "time"

type ContactSubmission struct {
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
