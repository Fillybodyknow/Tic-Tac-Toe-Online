package service

import (
	"errors"
	"tic-tac-toe-game/internal/model"
)

type GameLogicService struct {
	GameService   GameService
	PlayerService PlayerService
}

func (s *GameLogicService) MakeMove(RoomID string, PlayerID string, col int, row int) (*model.GameRoom, error) {
	var Symbol string
	GameRoom, err := s.GameService.FindGameRoomByID(RoomID)
	if err != nil {
		return nil, err
	}
	if GameRoom.PlayerX.PlayerID == PlayerID {
		Symbol = "X"
	} else if GameRoom.PlayerO.PlayerID == PlayerID {
		Symbol = "O"
	} else {
		return nil, errors.New("player not found in room")
	}

	if GameRoom.Turn == PlayerID {
		if GameRoom.Board[row][col] == "" {
			GameRoom.Board[row][col] = Symbol
			if _, IsWin := s.GameService.CheckWin(GameRoom.Board); IsWin {
				GameRoom.Winner = PlayerID
				P1, P2, _ := s.GameService.FindAllPlayersInRoom(RoomID)
				if GameRoom.Winner == P1.PlayerID {
					PlayerWinner, _ := s.PlayerService.FindPlayerByID(P1.PlayerID)
					PlayerLoser, _ := s.PlayerService.FindPlayerByID(P2.PlayerID)
					PlayerWinner.Win += 1
					PlayerLoser.Lose += 1
				} else {
					PlayerWinner, _ := s.PlayerService.FindPlayerByID(P2.PlayerID)
					PlayerLoser, _ := s.PlayerService.FindPlayerByID(P1.PlayerID)
					PlayerWinner.Win += 1
					PlayerLoser.Lose += 1
				}
				return GameRoom, nil
			}
			if s.GameService.CheckDraw(GameRoom.Board) {
				GameRoom.IsDraw = true
				P1, P2, _ := s.GameService.FindAllPlayersInRoom(RoomID)
				PlayerX, _ := s.PlayerService.FindPlayerByID(P1.PlayerID)
				PlayerO, _ := s.PlayerService.FindPlayerByID(P2.PlayerID)
				PlayerX.Draw += 1
				PlayerO.Draw += 1
				return GameRoom, nil
			}
			if GameRoom.Turn == GameRoom.PlayerX.PlayerID {
				GameRoom.Turn = GameRoom.PlayerO.PlayerID
			} else {
				GameRoom.Turn = GameRoom.PlayerX.PlayerID
			}
			return GameRoom, nil
		} else {
			return nil, errors.New("cell is not empty")
		}
	} else {
		return nil, errors.New("not your turn")
	}
}

func (s *GameService) CheckWin(board [3][3]string) (string, bool) {
	for i := 0; i < 3; i++ {
		if board[i][0] != "" && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0], true
		}
	}

	for i := 0; i < 3; i++ {
		if board[0][i] != "" && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			return board[0][i], true
		}
	}

	if board[0][0] != "" && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0], true
	}
	if board[0][2] != "" && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2], true
	}

	return "", false
}

func (s *GameService) CheckDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
