package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"tic-tac-toe-game/internal/hanlder"
	"tic-tac-toe-game/internal/router"
	"tic-tac-toe-game/internal/service"
	"tic-tac-toe-game/internal/socket"
)

func main() {

	socket.InitSocketServer(&hanlder.GameHandler{})
	server := socket.GetSocketServer()

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://wren-super-cobra.ngrok-free.app", "https://tic-tac-toe-online-tau.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	Handler := hanlder.NewPlayerHandler(service.PlayerService{})
	Router := router.NewRouter(Handler)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	r.StaticFS("/public", http.Dir("./public"))

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	API := r.Group("/api")
	{
		Router.PlayerRoute(API)
	}

	if err := r.Run("192.168.1.10:3000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
