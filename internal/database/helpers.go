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
func ScanRowUser(s scanner) (*models.User, error) {
	u := models.User{}
	var lastName, picture, phoneNumber, address, password sql.NullString
	var verifiedEmailToken, passwordResetToken sql.NullString

	err := s.Scan(
		&u.Id,
		&u.FirstName,
		&lastName,
		&u.Username,
		&u.Email,
		&password,
		&phoneNumber,
		&picture,
		&address,
		&u.IsActive,
		&u.VerifiedEmail,
		&verifiedEmailToken,
		&u.VerifiedEmailTokenExpiry,
		&passwordResetToken,
		&u.PasswordResetTokenExpiry,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if lastName.Valid {
		u.LastName = lastName.String
	}

	if picture.Valid {
		u.Picture = picture.String
	}

	if phoneNumber.Valid {
		u.PhoneNumber = phoneNumber.String
	}

	if address.Valid {
		u.Address = address.String
	}

	if password.Valid {
		u.Password = password.String
	}

	if verifiedEmailToken.Valid {
		u.VerifiedEmailToken = verifiedEmailToken.String
	}

	if passwordResetToken.Valid {
		u.PasswordResetToken = passwordResetToken.String
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// ScanRowRole scans a row into a Role struct
func ScanRowRole(s scanner) (*models.Role, error) {
	r := models.Role{}

	err := s.Scan(
		&r.Id,
		&r.Name,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// ScanRowUserRole scans a row into a UserRole struct
func ScanRowUserRole(s scanner) (*models.UserRole, error) {
	ur := models.UserRole{}

	err := s.Scan(
		&ur.UserId,
		&ur.RoleId,
	)
	if err != nil {
		return nil, err
	}

	return &ur, nil
}
