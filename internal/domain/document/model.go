package document

import "time"

// Document - структура описывающая документ.
type Document struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	File      bool      `json:"file"`
	Public    bool      `json:"public"`
	Mime      string    `json:"mime"`
	OwnerID   string    `json:"owner_id"`
	Grant     []string  `json:"grant"`
	CreatedAt time.Time `json:"created"`
}
