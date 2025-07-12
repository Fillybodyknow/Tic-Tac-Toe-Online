package service

import (
	"errors"
	"fmt"
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/utility"
)

type GameService struct{}

var GameRooms = []model.GameRoom{}

func (s *GameService) CreateRoom(Room model.GameRoom) model.GameRoom {
	var RoomID string
	for {
		RoomID = utility.GenerateRoomID()
		if _, err := s.FindGameRoomByID(RoomID); err != nil {
			break
		}
	}
	Room.RoomID = RoomID
	GameRooms = append(GameRooms, Room)
	return Room
}

func (s *GameService) FindGameRoomByID(RoomID string) (*model.GameRoom, error) {
	for i := range GameRooms {
		if GameRooms[i].RoomID == RoomID {
			return &GameRooms[i], nil
		}
	}
	return nil, errors.New("room not found")
}

func (s *GameService) ShowAllRooms() []model.GameRoom {
	return GameRooms
}

func (s *GameService) JoinGame(RoomID string, Player model.Player) error {
	GameRoom, err := s.FindGameRoomByID(RoomID)
	if err != nil {
		fmt.Println("Room not found")
		return errors.New("room not found")
	}
	if GameRoom.PlayerX.PlayerID == "" {
		GameRoom.PlayerX = Player
	} else if GameRoom.PlayerO.PlayerID == "" {
		GameRoom.PlayerO = Player
	} else {
		return errors.New("room is full")
	}
	return nil
}

func (s *GameService) FindPlayerByIDInGameRoom(RoomID string, PlayerID string) (model.Player, error) {
	var HavePlayer model.Player
	for _, p := range GameRooms {
		if p.RoomID == RoomID {
			if p.PlayerX.PlayerID == PlayerID {
				HavePlayer = p.PlayerX
				fmt.Println(HavePlayer)
				return HavePlayer, nil
			} else if p.PlayerO.PlayerID == PlayerID {
				HavePlayer = p.PlayerO
				fmt.Println(HavePlayer)
				return HavePlayer, nil
			}
		}
	}
	fmt.Println("Player not found in room")
	return model.Player{}, errors.New("player not found")
}

func (s *GameService) StartGame(RoomID string) error {
	GameRoom, err := s.FindGameRoomByID(RoomID)
	if err != nil {
		return err
	}
	if GameRoom.PlayerX.PlayerID == "" || GameRoom.PlayerO.PlayerID == "" {
		return errors.New("game room is not ready to start")
	}
	for i := range GameRooms {
		if GameRooms[i].RoomID == RoomID {
			if GameRooms[i].Turn == "" {
				GameRooms[i].Turn = GameRoom.PlayerX.PlayerID
			}
		}
	}
	return nil
}

func (s *GameService) LeaveGameRoom(RoomID string, PlayerID string) error {
	GameRoom, err := s.FindGameRoomByID(RoomID)
	if err != nil {
		return err
	}
	if GameRoom.PlayerX.PlayerID == PlayerID {
		GameRoom.PlayerX = model.Player{}
	} else if GameRoom.PlayerO.PlayerID == PlayerID {
		GameRoom.PlayerO = model.Player{}
	} else {
		return errors.New("player not found in room")
	}
	return nil
}

func (s *GameService) FindAllPlayersInRoom(RoomID string) (model.Player, model.Player, error) {
	GameRoom, err := s.FindGameRoomByID(RoomID)
	if err != nil {
		return model.Player{}, model.Player{}, err
	}
	return GameRoom.PlayerX, GameRoom.PlayerO, nil
}

func (s *GameService) DeleteRoom(RoomID string) error {
	for i := range GameRooms {
		if GameRooms[i].RoomID == RoomID {
			GameRooms = append(GameRooms[:i], GameRooms[i+1:]...)
			return nil
		}
	}
	return errors.New("room not found")
}

func (s *GameService) ResetBoard(RoomID string) error {
	for i := range GameRooms {
		if GameRooms[i].RoomID == RoomID {
			GameRooms[i].Board = [3][3]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}
			GameRooms[i].Special_PawnX = map[string]int{
				"X|medium|2": 2,
				"X|large|3":  1,
			}
			GameRooms[i].Special_PawnO = map[string]int{
				"O|medium|2": 2,
				"O|large|3":  1,
			}
			return nil
		}
	}
	return errors.New("room not found")
}

func (s *GameService) ResetTurn(RoomID string) error {
	for i := range GameRooms {
		if GameRooms[i].RoomID == RoomID {
			GameRooms[i].Turn = ""
			return nil
		}
	}
	return errors.New("room not found")
}
