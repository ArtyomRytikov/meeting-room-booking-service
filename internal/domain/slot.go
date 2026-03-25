package domain

import "time"

type Slot struct {
	ID     string `json:"id"`
	RoomID string `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

type SlotDetails struct {
	SlotID  string
	RoomID  string
	StartAt time.Time
	EndAt   time.Time
}
