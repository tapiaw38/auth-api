package database

import (
	"database/sql"

	"github.com/tapiaw38/auth-api/internal/models"
)

// scanner is the interface that wraps the basic Scan method.
type scanner interface {
	Scan(dest ...interface{}) error
}

// ScanRowUser scans a row into a User struct
func ScanRowUser(s scanner) (models.User, error) {
	u := models.User{}
	var lastName, picture, phoneNumber, address sql.NullString
	var isActive, verifiedEmail sql.NullBool

	err := s.Scan(
		&u.Id,
		&u.FirstName,
		&lastName,
		&u.Username,
		&u.Email,
		&u.Password,
		&phoneNumber,
		&picture,
		&address,
		&isActive,
		&verifiedEmail,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	u.LastName = lastName.String
	u.Picture = picture.String
	u.PhoneNumber = phoneNumber.String
	u.Address = address.String

	return u, nil
}

func ScanRowUserResponse(s scanner) (models.UserResponse, error) {
	u := models.UserResponse{}
	var lastName, picture, phoneNumber, address sql.NullString
	var isActive, verifiedEmail sql.NullBool
	var createdAt, updatedAt sql.NullTime

	err := s.Scan(
		&u.Id,
		&u.FirstName,
		&lastName,
		&u.Username,
		&u.Email,
		&phoneNumber,
		&picture,
		&address,
		&isActive,
		&verifiedEmail,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return u, err
	}

	u.LastName = lastName.String
	u.Picture = picture.String
	u.PhoneNumber = phoneNumber.String
	u.Address = address.String
	u.IsActive = isActive.Bool
	u.VerifiedEmail = verifiedEmail.Bool
	u.CreatedAt = createdAt.Time
	u.UpdatedAt = updatedAt.Time

	return u, nil
}

// ScanRowRole scans a row into a Role struct
func ScanRowRole(s scanner) (models.Role, error) {
	r := models.Role{}

	err := s.Scan(
		&r.Id,
		&r.Name,
	)
	if err != nil {
		return r, err
	}

	return r, nil
}

// ScanRowUserRole scans a row into a UserRole struct
func ScanRowUserRole(s scanner) (models.UserRole, error) {
	ur := models.UserRole{}

	err := s.Scan(
		&ur.UserId,
		&ur.RoleId,
	)
	if err != nil {
		return ur, err
	}

	return ur, nil
}
