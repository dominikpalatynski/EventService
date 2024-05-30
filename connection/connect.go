package connection

import (
	"net/http"

	"github.com/dominikpalatynski/EventService/storage"
	"github.com/gin-gonic/gin"
)

type Event struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

type APIServer struct {
	router *gin.Engine
	storage *storage.StorePostgres
}

func NewAPIServer(s *storage.StorePostgres) *APIServer{
	r := gin.Default()
	server := &APIServer{
		router: r,
		storage: s,
	}
	return server
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
		Name : "elo",
		Content:  "test",
	}
	c.JSON(http.StatusOK, event)
}