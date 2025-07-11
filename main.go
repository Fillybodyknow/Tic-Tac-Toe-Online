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
	r := gin.Default()

	r.Use(cors.Default())

	socket.InitSocketServer(&hanlder.GameHandler{})

	server := socket.GetSocketServer()

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

	if err := r.Run(":3000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
