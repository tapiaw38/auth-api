package handlers

import (
	"github.com/gin-gonic/gin"
)

const (
	Error   = "error"
	Message = "message"
)

// Response is the struct that will be sent to the client
type Response struct {
	MessageType string      `json:"message_type"`
	Message     string      `json:"message"`
	Response    interface{} `json:"response"`
}

// NewResponse creates a new Response struct
func NewResponse(messageType string, message string, response interface{}) Response {
	return Response{
		MessageType: messageType,
		Message:     message,
		Response:    response,
	}
}

// ResponseWithJson sends a json response to the client
func ResponseWithJson(c *gin.Context, statusCode int, response Response) {
	c.JSON(statusCode, response)
}
