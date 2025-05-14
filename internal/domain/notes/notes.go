package notes

type Note struct {
	UID         string `json:"uid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}