package notes

import "time"

type Note struct {
	NID         string `json:"nid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	Created_at  time.Time `json:"created_at"`
	UID         string `json:"uid"`
}
