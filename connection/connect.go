package connection

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dominikpalatynski/EventService/storage"
	"github.com/dominikpalatynski/EventService/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	util.LoadEnv()

	s.router.Run(":"+ os.Getenv("PORT"))
}

func (s *APIServer) registerRoutes() {
	s.router.GET("/events", s.getEvents)
	s.router.POST("/events", s.addEvent)
	s.router.PATCH("/events/:id", s.updateEvent)
	s.router.DELETE("events/:id", s.deleteEvent)
}

func (s *APIServer) getEvents(c *gin.Context) {
	events, err := s.storage.GetEvents()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, events)
}

func (s *APIServer) addEvent(c *gin.Context) {
	event := new(storage.Event)
	
	userId, ok := cookieReader("UserId", c)
	
	if ok != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ok.Error()})
		return
	}
	event.UserId = userId

	if ok := getEventFromAPI(c, event); ok != nil {
		return
	}

	err := s.storage.AddEvent(event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (s *APIServer) updateEvent(c *gin.Context) {
	id := c.Param("id")
	eventId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedData map[string]interface{}

	if ok := getDataToUpdate(c, &updatedData); ok != nil {
		return
	}

	if content, ok := updatedData["content"]; ok {
        contentBytes, err := json.Marshal(content)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content"})
            return
        }
        updatedData["content"] = contentBytes
    }

	statusOk := s.storage.UpdateById(eventId, updatedData)

	if statusOk != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": statusOk.Error()})
		return
	}

	c.JSON(http.StatusOK, "event updated succesfully")
}

func (s *APIServer) deleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusOk := s.storage.DeleteById(eventId)

	if statusOk != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": statusOk.Error()})
		return
	}

	c.JSON(http.StatusOK, "event updated succesfully")
}