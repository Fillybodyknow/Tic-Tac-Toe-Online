package model

import "time"

type Player struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Win        int    `json:"win"`
	Lose       int    `json:"lose"`
	Draw       int    `json:"draw"`
}

type GameRoom struct {
	RoomID        string         `json:"room_id"`
	PlayerX       Player         `json:"player_x"`
	PlayerO       Player         `json:"player_o"`
	Board         [3][3]string   `json:"board"`
	Special_PawnX map[string]int `json:"special_pawn_x"`
	Special_PawnO map[string]int `json:"special_pawn_o"`
	Turn          string         `json:"turn"`
	Winner        string         `json:"winner"`
	IsDraw        bool           `json:"is_draw"`
	CreatedAt     time.Time      `json:"created_at"`
}

type PlayerInput struct {
	PlayerName string `json:"username"`
}
