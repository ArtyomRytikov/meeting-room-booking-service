package domain

type Slot struct {
	ID     string `json:"id"`
	RoomID string `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}
