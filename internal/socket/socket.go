package socket

import (
	"fmt"
	"log"
	"tic-tac-toe-game/internal/hanlder"

	socketio "github.com/googollee/go-socket.io"
)

var (
	server *socketio.Server
)

type MakeMovePayload struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

func InitSocketServer(gameHandler *hanlder.GameHandler) {
	server = socketio.NewServer(nil)
	GameRoom := "/game-room"

	server.OnConnect("/", func(Conn socketio.Conn) error {
		Client := Conn.ID()
		fmt.Println("Have User Connection :" + Client)
		Conn.Emit("connect-successfuly", "Connected Successfuly... :"+Client)
		AllRooms := gameHandler.ShowAllGameRoom()
		for _, room := range AllRooms {
			if room.PlayerX.PlayerID == "" || room.PlayerO.PlayerID == "" {
				gameHandler.ResetGame(room.RoomID)
			}
		}
		server.BroadcastToNamespace("/", "connection", AllRooms)
		return nil
	})

	server.OnEvent("/", "create-room", func(Conn socketio.Conn) {
		u := Conn.URL()
		playerID := u.Query().Get("player_id")
		fmt.Println("This Client: " + playerID + " - Create Room")
		NewRoom := gameHandler.CreateGameRoom()
		data := map[string]interface{}{
			"room":      NewRoom,
			"player_id": playerID,
		}
		server.BroadcastToNamespace("/", "create-room-successfuly", data)
	})

	server.OnConnect(GameRoom, func(Conn socketio.Conn) error {
		Client := Conn.ID()
		u := Conn.URL()
		playerID := u.Query().Get("player_id")
		roomID := u.Query().Get("room_id")
		Conn.SetContext(map[string]string{
			"player_id": playerID,
			"room_id":   roomID,
		})
		if playerID == "" || roomID == "" {
			Conn.Emit("join-room-failed", "Missing player_id or room_id")
			return fmt.Errorf("invalid connection query")
		}

		Conn.Join(roomID)
		err := gameHandler.JoinGameRoom(roomID, playerID)
		if err != nil {
			Conn.Leave(roomID)
			Conn.Emit("join-room-failed", err.Error())
			return err
		}
		GameIsReady := gameHandler.StartGame(roomID)
		log.Println("GameIsReady: ", GameIsReady)
		if GameIsReady {
			server.BroadcastToRoom(GameRoom, roomID, "game-ready", "Game Is Ready...")
		} else {
			server.BroadcastToRoom(GameRoom, roomID, "game-not-ready", "Waiting For Other Player...")
		}
		fmt.Println("Have User : " + playerID + " Connection in Game Room :" + roomID)
		Conn.Emit("connect-successfuly", "Connected Successfuly... :"+Client)
		AllRooms := gameHandler.ShowAllGameRoom()
		server.BroadcastToNamespace("/", "connection", AllRooms)
		return nil
	})

	server.OnDisconnect(GameRoom, func(Conn socketio.Conn, reason string) {
		ctx := Conn.Context().(map[string]string)
		roomID := ctx["room_id"]
		playerID := ctx["player_id"]

		gameHandler.LeaveRoom(roomID, playerID)

		server.BroadcastToRoom(GameRoom, roomID, "game-not-ready", "Waiting For Other Player...")
		AllRooms := gameHandler.ShowAllGameRoom()
		server.BroadcastToNamespace("/", "connection", AllRooms)
	})

	server.OnEvent(GameRoom, "make-move", func(Conn socketio.Conn, payload MakeMovePayload) {
		ctx := Conn.Context().(map[string]string)
		roomID := ctx["room_id"]
		playerID := ctx["player_id"]

		updateRoom, err := gameHandler.MakeMove(roomID, playerID, payload.Col, payload.Row)
		if err != nil {
			Conn.Emit("make-move-failed", err.Error())
			return
		}

		server.BroadcastToRoom(GameRoom, roomID, "update-board", updateRoom)

		if updateRoom.Winner != "" {
			server.BroadcastToRoom(GameRoom, roomID, "game-winner", updateRoom.Winner)
			return
		}

		if updateRoom.IsDraw {
			server.BroadcastToRoom(GameRoom, roomID, "game-draw", "The game is a draw.")
		}
	})

	server.OnEvent(GameRoom, "request-game-state", func(conn socketio.Conn) {
		ctx := conn.Context().(map[string]string)
		roomID := ctx["room_id"]

		room, err := gameHandler.GameService.FindGameRoomByID(roomID)
		if err != nil {
			conn.Emit("make-move-failed", "Room not found")
			return
		}

		conn.Emit("update-board", room)
	})

}

func GetSocketServer() *socketio.Server {
	return server
}
