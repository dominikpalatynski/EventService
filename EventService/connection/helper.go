package connection

import (
	"net/http"

	"github.com/dominikpalatynski/EventService/types"
	"github.com/gin-gonic/gin"
)

func cookieReader(variable string, ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(variable)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func getEventFromAPI(c *gin.Context, event *types.Event) error{
	if err := c.ShouldBindJSON(event); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
	}
	return nil
}

func getDataToUpdate(c *gin.Context, event *map[string]interface{}) error{
	if err := c.ShouldBindJSON(event); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
	}
	return nil
}