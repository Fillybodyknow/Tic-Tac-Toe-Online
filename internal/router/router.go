package router

import (
	"tic-tac-toe-game/internal/hanlder"

	"github.com/gin-gonic/gin"
)

type Router struct {
	PlayerHandler *hanlder.PlayerHandler
}

func NewRouter(playerHandler *hanlder.PlayerHandler) *Router {
	return &Router{PlayerHandler: playerHandler}
}

func (R *Router) PlayerRoute(r *gin.RouterGroup) {
	r.POST("/player", R.PlayerHandler.CreatePlayer)
	r.GET("/player/:player_id", R.PlayerHandler.ShowPlayerData)
}
