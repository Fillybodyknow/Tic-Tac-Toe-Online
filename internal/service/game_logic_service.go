package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tic-tac-toe-game/internal/model"
)

type GameLogicService struct {
	GameService   GameService
	PlayerService PlayerService
}

func (s *GameLogicService) MakeMove(RoomID, Special_Pawn, PlayerID string, col, row int) (*model.GameRoom, error) {
	GameRoom, err := s.GameService.FindGameRoomByID(RoomID)
	if err != nil {
		return nil, err
	}

	Symbol, err := s.getPlayerSymbol(GameRoom, PlayerID, Special_Pawn)
	if err != nil {
		return nil, err
	}

	if GameRoom.Turn != PlayerID {
		return nil, errors.New("not your turn")
	}

	PriorityYourPawn := s.extractPawnPriority(Symbol)

	cellContent := GameRoom.Board[row][col]
	if cellContent == "" {
		GameRoom.Board[row][col] = Symbol
		if len(Symbol) > 1 {
			if strings.HasPrefix(Symbol, "X|") && GameRoom.PlayerX.PlayerID == PlayerID {
				if GameRoom.Special_PawnX[Special_Pawn] > 0 {
					GameRoom.Special_PawnX[Special_Pawn] -= 1
					fmt.Println("Special_PawnX:", GameRoom.Special_PawnX)
				} else {
					return nil, errors.New("your selected pawn is not available")
				}
			} else if strings.HasPrefix(Symbol, "O|") && GameRoom.PlayerO.PlayerID == PlayerID {
				if GameRoom.Special_PawnO[Special_Pawn] > 0 {
					GameRoom.Special_PawnO[Special_Pawn] -= 1
					fmt.Println("Special_PawnO:", GameRoom.Special_PawnO)
				} else {
					return nil, errors.New("your selected pawn is not available")
				}
			} else {
				return nil, errors.New("invalid special pawn usage")
			}
		}
		return s.handleWinOrDraw(GameRoom, PlayerID)
	}

	OpponentSymbol, _ := s.GetSymbol_OnRowCol(RoomID, PlayerID, col, row)
	if s.isYourOwnPawn(Symbol, OpponentSymbol) {
		return nil, errors.New("this is your pawn")
	}

	PriorityOpponentPawn := s.extractPawnPriority(OpponentSymbol)
	if PriorityYourPawn > PriorityOpponentPawn {
		GameRoom.Board[row][col] = Symbol
		if len(Symbol) > 1 {
			if strings.HasPrefix(Symbol, "X|") && GameRoom.PlayerX.PlayerID == PlayerID {
				if GameRoom.Special_PawnX[Special_Pawn] > 0 {
					GameRoom.Special_PawnX[Special_Pawn] -= 1
					fmt.Println("Special_PawnX:", GameRoom.Special_PawnX)
				} else {
					return nil, errors.New("your selected pawn is not available")
				}
			} else if strings.HasPrefix(Symbol, "O|") && GameRoom.PlayerO.PlayerID == PlayerID {
				if GameRoom.Special_PawnO[Special_Pawn] > 0 {
					GameRoom.Special_PawnO[Special_Pawn] -= 1
					fmt.Println("Special_PawnO:", GameRoom.Special_PawnO)
				} else {
					return nil, errors.New("your selected pawn is not available")
				}
			} else {
				return nil, errors.New("invalid special pawn usage")
			}
		}
		return s.handleWinOrDraw(GameRoom, PlayerID)
	}

	return nil, errors.New("your pawn is blocked by opponent's pawn")
}

func (s *GameLogicService) GetSymbol_OnRowCol(RoomID string, PlayerID string, col int, row int) (string, error) {
	GameRoom, err := s.GameService.FindGameRoomByID(RoomID)
	if err != nil {
		return "", err
	}
	fmt.Println("Symbol on cell:", GameRoom.Board[row][col])
	return GameRoom.Board[row][col], nil
}

func (s *GameLogicService) getPlayerSymbol(GameRoom *model.GameRoom, PlayerID, Special_Pawn string) (string, error) {
	if GameRoom.PlayerX.PlayerID == PlayerID {
		if Special_Pawn != "" {
			fmt.Println("Your Special Pawn:", Special_Pawn)
			return Special_Pawn, nil
		}
		fmt.Println("Your Symbol: X")
		return "X", nil
	} else if GameRoom.PlayerO.PlayerID == PlayerID {
		if Special_Pawn != "" {
			fmt.Println("Your Special Pawn:", Special_Pawn)
			return Special_Pawn, nil
		}
		fmt.Println("Your Symbol: O")
		return "O", nil
	}
	return "", errors.New("player not found in room")
}

func (s *GameLogicService) extractPawnPriority(symbol string) int {
	parts := strings.Split(symbol, "|")
	if len(parts) >= 3 {
		if priority, err := strconv.Atoi(parts[2]); err == nil {
			fmt.Println("Priority:", priority)
			return priority
		}
	}
	fmt.Println("Priority: 1")
	return 1
}

func (s *GameLogicService) isYourOwnPawn(yourSymbol, opponentSymbol string) bool {
	yourParts := strings.Split(yourSymbol, "|")
	opponentParts := strings.Split(opponentSymbol, "|")
	if len(yourParts) > 0 && len(opponentParts) > 0 {
		fmt.Println("isYourOwnPawn:", yourParts[0] == opponentParts[0])
		return yourParts[0] == opponentParts[0]
	}
	fmt.Println("isYourOwnPawn: false")
	return false
}

func (s *GameLogicService) switchTurn(GameRoom *model.GameRoom) {
	if GameRoom.Turn == GameRoom.PlayerX.PlayerID {
		GameRoom.Turn = GameRoom.PlayerO.PlayerID
	} else {
		GameRoom.Turn = GameRoom.PlayerX.PlayerID
	}
}

func (s *GameLogicService) updateWinLoseScore(GameRoom *model.GameRoom, WinnerID string) {
	P1, P2, _ := s.GameService.FindAllPlayersInRoom(GameRoom.RoomID)
	PlayerWinnerID := P1.PlayerID
	PlayerLoserID := P2.PlayerID
	if WinnerID != P1.PlayerID {
		PlayerWinnerID = P2.PlayerID
		PlayerLoserID = P1.PlayerID
	}

	Winner, _ := s.PlayerService.FindPlayerByID(PlayerWinnerID)
	Loser, _ := s.PlayerService.FindPlayerByID(PlayerLoserID)
	Winner.Win++
	Loser.Lose++
}

func (s *GameLogicService) updateDrawScore(GameRoom *model.GameRoom) {
	P1, P2, _ := s.GameService.FindAllPlayersInRoom(GameRoom.RoomID)
	PlayerX, _ := s.PlayerService.FindPlayerByID(P1.PlayerID)
	PlayerO, _ := s.PlayerService.FindPlayerByID(P2.PlayerID)
	PlayerX.Draw++
	PlayerO.Draw++
}

func (s *GameLogicService) handleWinOrDraw(GameRoom *model.GameRoom, PlayerID string) (*model.GameRoom, error) {
	if _, IsWin := s.CheckWin(GameRoom.Board); IsWin {
		GameRoom.Winner = PlayerID
		s.updateWinLoseScore(GameRoom, PlayerID)
		return GameRoom, nil
	}
	if s.CheckDraw(GameRoom.Board) {
		GameRoom.IsDraw = true
		s.updateDrawScore(GameRoom)
		return GameRoom, nil
	}
	s.switchTurn(GameRoom)
	return GameRoom, nil
}

func (s *GameLogicService) GetNormalSymbol(Symbol string) string {
	if len(Symbol) > 1 {
		Symbol := strings.Split(Symbol, "|")[0]
		return Symbol
	}
	return Symbol
}

func (s *GameLogicService) CheckWin(board [3][3]string) (string, bool) {
	for i := 0; i < 3; i++ {
		if board[i][0] != "" && s.GetNormalSymbol(board[i][0]) == s.GetNormalSymbol(board[i][1]) && s.GetNormalSymbol(board[i][1]) == s.GetNormalSymbol(board[i][2]) {
			return board[i][0], true
		}
	}

	for i := 0; i < 3; i++ {
		if board[0][i] != "" && s.GetNormalSymbol(board[0][i]) == s.GetNormalSymbol(board[1][i]) && s.GetNormalSymbol(board[1][i]) == s.GetNormalSymbol(board[2][i]) {
			return board[0][i], true
		}
	}

	if board[0][0] != "" && s.GetNormalSymbol(board[0][0]) == s.GetNormalSymbol(board[1][1]) && s.GetNormalSymbol(board[1][1]) == s.GetNormalSymbol(board[2][2]) {
		return board[0][0], true
	}

	if board[0][2] != "" && s.GetNormalSymbol(board[0][2]) == s.GetNormalSymbol(board[1][1]) && s.GetNormalSymbol(board[1][1]) == s.GetNormalSymbol(board[2][0]) {
		return board[0][2], true
	}

	return "", false
}

func (s *GameLogicService) CheckDraw(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
