package hanlder

import (
	"errors"
	"fmt"
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/service"
	"time"
)

type GameHandler struct {
	GameService      service.GameService
	PlayerService    service.PlayerService
	GameLogicService service.GameLogicService
}

func NewGameHanlder(gameService service.GameService, playerService service.PlayerService, gameLogicService service.GameLogicService) *GameHandler {
	return &GameHandler{GameService: gameService, PlayerService: playerService, GameLogicService: gameLogicService}
}

func (H *GameHandler) CreateGameRoom() model.GameRoom {
	GameRoom := model.GameRoom{
		RoomID: "",
		Special_PawnX: map[string]int{
			"x|medium": 3,
			"x|large":  2,
		},
		Special_PawnO: map[string]int{
			"o|medium": 3,
			"o|large":  2,
		},
		PlayerX:   model.Player{},
		PlayerO:   model.Player{},
		Board:     [3][3]string{},
		Turn:      "",
		Winner:    "",
		IsDraw:    false,
		CreatedAt: time.Now(),
	}

	CreatedRoom := H.GameService.CreateRoom(GameRoom)
	return CreatedRoom

}

func (H *GameHandler) ShowAllGameRoom() []model.GameRoom {
	return H.GameService.ShowAllRooms()
}

func (H *GameHandler) JoinGameRoom(RoomID string, PlayerID string) error {
	if _, err := H.GameService.FindPlayerByIDInGameRoom(RoomID, PlayerID); err == nil {
		return errors.New("player already in room")
	}

	Player, err := H.PlayerService.FindPlayerByID(PlayerID)
	if err != nil {
		fmt.Println("Not Found Player")
		return err
	}
	if err := H.GameService.JoinGame(RoomID, *Player); err != nil {
		fmt.Println("Not Found Room")
		return err
	}
	H.GameService.FindGameRoomByID(RoomID)
	return nil
}

func (H *GameHandler) LeaveRoom(RoomID string, PlayerID string) bool {
	if err := H.GameService.LeaveGameRoom(RoomID, PlayerID); err != nil {
		fmt.Println("LeaveRoom error:", err)
		return false
	}
	PlayerX, PlayerO, err := H.GameService.FindAllPlayersInRoom(RoomID)
	if err != nil {
		fmt.Println("LeaveRoom error:", err)
		return false
	}
	if PlayerX.PlayerID == "" && PlayerO.PlayerID == "" {
		if err := H.GameService.DeleteRoom(RoomID); err != nil {
			fmt.Println("LeaveRoom error:", err)
			return false
		}
	}
	return true
}

func (H *GameHandler) StartGame(RoomID string) bool {
	if err := H.GameService.StartGame(RoomID); err != nil {
		return false
	}
	return true
}

func (H *GameHandler) MakeMove(RoomID string, PlayerID string, col int, row int) (*model.GameRoom, error) {
	updateRoom, err := H.GameLogicService.MakeMove(RoomID, PlayerID, col, row)
	if err != nil {
		return nil, err
	}
	return updateRoom, nil
}

func (H *GameHandler) ResetGame(RoomID string) error {
	if err := H.GameService.ResetBoard(RoomID); err != nil {
		return err
	}
	if err := H.GameService.ResetTurn(RoomID); err != nil {
		return err
	}
	return nil
}
