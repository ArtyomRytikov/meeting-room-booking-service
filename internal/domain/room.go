package domain

type Room struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
}
