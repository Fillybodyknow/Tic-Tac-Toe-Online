package hanlder

import (
	"net/http"
	"tic-tac-toe-game/internal/model"
	"tic-tac-toe-game/internal/service"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	PlayerService service.PlayerService
}

func NewPlayerHandler(playerService service.PlayerService) *PlayerHandler {
	return &PlayerHandler{PlayerService: playerService}
}

func (H *PlayerHandler) CreatePlayer(c *gin.Context) {
	var PlayerInput = model.PlayerInput{}

	if err := c.ShouldBindJSON(&PlayerInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if FoundPlayer, err := H.PlayerService.FindPlayerByName(PlayerInput.PlayerName); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"player": FoundPlayer})
		return
	}
	Player := model.Player{
		PlayerName: PlayerInput.PlayerName,
		Win:        0,
		Lose:       0,
		Draw:       0,
	}
	PlayerCreated := H.PlayerService.InsertPlayer(Player)
	c.JSON(http.StatusCreated, gin.H{"player": PlayerCreated})
}

func (H *PlayerHandler) ShowPlayerData(c *gin.Context) {
	userID := c.Param("player_id")

	PlayerData, err := H.PlayerService.FindPlayerByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"player": PlayerData})
}
