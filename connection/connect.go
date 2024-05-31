package connection

import (
	"log"
	"net/http"
	"os"

	"github.com/dominikpalatynski/EventService/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


type APIServer struct {
	router *gin.Engine
	storage *storage.MongoDbStorage
}

func NewAPIServer(s *storage.MongoDbStorage) *APIServer{
	r := gin.Default()
	server := &APIServer{
		router: r,
		storage: s,
	}
	return server
}

func (s *APIServer) Run() {

	s.registerRoutes()

	if err := godotenv.Load(".env"); err !=nil {
		log.Fatal("Error loading .env")
	}

	s.router.Run(":"+ os.Getenv("PORT"))
}

func (s *APIServer) registerRoutes() {
	s.router.GET("/events", s.getEvents)
	s.router.POST("/events", s.addEvent)
}

func (s *APIServer) getEvents(c *gin.Context) {
	events, err := s.storage.GetEvents()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, events)
}

func (s *APIServer) addEvent(c *gin.Context) {
	title := new(storage.Event)

	if err := c.ShouldBindJSON(title); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
	}

	event, err := s.storage.AddEvent(title.Title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, event)
}