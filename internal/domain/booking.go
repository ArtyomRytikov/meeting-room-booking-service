package domain

type Booking struct {
	ID        string `json:"id"`
	SlotID    string `json:"slotId"`
	RoomID    string `json:"roomId"`
	UserID    string `json:"userId"`
	Status    string `json:"status"`
	Start     string `json:"start"`
	End       string `json:"end"`
	CreatedAt string `json:"createdAt"`
}
