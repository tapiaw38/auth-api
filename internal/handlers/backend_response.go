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

// handleError sends an error response
func HandleError(c *gin.Context, code int, err error) {
	response := NewResponse(Error, err.Error(), nil)
	ResponseWithJson(c, code, response)
}

// HandleSuccess sends a success response
func HandleSuccess(c *gin.Context, code int, message string, data interface{}) {
	response := NewResponse(Message, message, data)
	ResponseWithJson(c, code, response)
}
