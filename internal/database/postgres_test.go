package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	DATABASE_URL = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
)

func TestNewPostgresRepository(t *testing.T) {
	t.Run("should return a postgres repository", func(t *testing.T) {
		repository, err := NewPostresRepository(DATABASE_URL)
		assert.Nil(t, err)
		assert.NotNil(t, repository)
	})
}

func TestGetRelativePathToMigrationsDirectory(t *testing.T) {
	t.Run("should return the relative path to migrations directory", func(t *testing.T) {
		path, err := GetRelativePathToMigrationsDirectory()
		assert.Nil(t, err)
		assert.NotNil(t, path)
		assert.Equal(t, "file://migrations", path)
	})
}

// func TestMakemigration(t *testing.T) {
// 	t.Run("should make the database migrations", func(t *testing.T) {
// 		db, mock, err := sqlmock.New()
// 		assert.Nil(t, err)
// 		defer db.Close()

// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS roles").WillReturnResult(sqlmock.NewResult(1, 1))
// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS user_roles").WillReturnResult(sqlmock.NewResult(1, 1))

// 		repository := &PostgresRepository{db}
// 		err = repository.Makemigration(DATABASE_URL)
// 		assert.Nil(t, err)
// 		if err != nil {
// 			t.Errorf("error was not expected while making migrations: %s", err)
// 		}
// 	})
// }

// func TestGetDB(t *testing.T) {
// 	t.Run("should return the database instance", func(t *testing.T) {
// 		db, mock, err := sqlmock.New()
// 		assert.Nil(t, err)
// 		defer db.Close()

// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))
// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS roles").WillReturnResult(sqlmock.NewResult(1, 1))
// 		mock.ExpectExec("CREATE TABLE IF NOT EXISTS user_roles").WillReturnResult(sqlmock.NewResult(1, 1))

// 		repository := &PostgresRepository{db}
// 		err = repository.Makemigration(DATABASE_URL)
// 		assert.Nil(t, err)

// 		dbInstance := repository.db
// 		assert.NotNil(t, dbInstance)
// 	})
// }
