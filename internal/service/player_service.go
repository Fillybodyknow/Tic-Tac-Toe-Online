package service

import (
	"errors"
	"fmt"
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/utility"
)

type PlayerService struct{}

var PlayerDatas = []model.Player{}

func (s *PlayerService) InsertPlayer(Player model.Player) model.Player {
	userID := utility.GenerateUserID()
	for {
		if _, err := s.FindPlayerByID(userID); err != nil {
			break
		}
	}
	Player.PlayerID = userID
	PlayerDatas = append(PlayerDatas, Player)
	return Player
}

func (s *PlayerService) FindPlayerByID(userID string) (*model.Player, error) {
	for i := range PlayerDatas {
		if PlayerDatas[i].PlayerID == userID {
			return &PlayerDatas[i], nil
		}
	}
	fmt.Println("Player not found")
	return nil, errors.New("player not found")
}

func (s *PlayerService) FindPlayerByName(userName string) (model.Player, error) {
	var HavePlayer model.Player
	for _, p := range PlayerDatas {
		if p.PlayerName == userName {
			HavePlayer = p
			return HavePlayer, nil
		}
	}
	return model.Player{}, errors.New("player not found")
}

func (s *PlayerService) GetAllPlayers() []model.Player {
	return PlayerDatas
}

func (s *PlayerService) UpdatePlayer(Player model.Player) model.Player {
	for i, p := range PlayerDatas {
		if p.PlayerID == Player.PlayerID {
			PlayerDatas[i] = Player
			return Player
		}
	}
	return model.Player{}
}
