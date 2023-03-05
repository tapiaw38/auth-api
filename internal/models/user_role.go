package models

// UserRole is the model for the user_role table
type UserRole struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}
