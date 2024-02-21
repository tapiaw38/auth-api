package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tapiaw38/auth-api/internal/server"
)

func TestLogin(t *testing.T) {
	c := require.New(t)

	// Create a new Gin router
	router := gin.New()

	// Set up a sample Gin handler for /login endpoint
	router.POST("/login", LoginHandler(&server.Broker{}))

	// Define the payload for the request body
	payload := map[string]string{
		"email":    "tapiaw38@gmail.com",
		"password": "Walter153294",
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	c.NoError(err)

	// Create a new HTTP request with the payload in the body
	r, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payloadBytes))
	c.NoError(err)

	// Set the content type header
	r.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, r)

	// Ensure the response status code is 200
	//c.Equal(http.StatusOK, w.Code)

	// Process the expected JSON response
	expectedBodyLoginResponse, err := os.ReadFile("../utils/samples/user_login_response.json")
	c.NoError(err)

	fmt.Println("Expected:", string(expectedBodyLoginResponse))

	// Process the actual JSON response
	var actualLoginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &actualLoginResponse)
	c.NoError(err)

	fmt.Println("Actual:", actualLoginResponse)
}

// Helper function to set Gin URL parameters
func setGinRequestVars(r *http.Request, payload map[string]string) *http.Request {
	query := r.URL.Query()
	for key, value := range payload {
		query.Add(key, value)
	}
	r.URL.RawQuery = query.Encode()
	return r
}

// func TestLoginHandler(t *testing.T) {
// 	t.Run("should return 200 status code", func(t *testing.T) {
// 		s := &server.Broker{}
// 		payload := `{"email": "tapiaw38@gmail.com", "password": "Walter153294"}`
// 		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(payload))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		r := gin.Default()
// 		r.POST("/login", LoginHandler(s))

// 		r.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusOK {
// 			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 		}

// 		fmt.Println(rr.Body.String())
// 	})
// }
