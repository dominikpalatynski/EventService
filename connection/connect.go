package connection

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

type APIServer struct {
	router *gin.Engine
}

func NewAPIServer() *APIServer{
	r := gin.Default()
	s := &APIServer{
		router: r,
	}
	return s
}

func (s *APIServer) Run() {

	s.registerRoutes()

	s.router.Run(":8080")
}

func (s *APIServer) registerRoutes() {
	s.router.GET("/events", s.getEvents)
}

func (s *APIServer) getEvents(c *gin.Context) {
	event := Event{
		Name : "test",
		Content:  "contents",
	}
	c.JSON(http.StatusOK, event)
}