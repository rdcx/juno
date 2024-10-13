package mysql

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"juno/pkg/api/user"
	"juno/pkg/api/user/migration/mysql"
	"testing"

	"github.com/google/uuid"
)

func longDummyString(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += "a"
	}
	return s
}

func randomEmail() string {
	return uuid.New().String() + "@example.com"
}

func testUserMatches(t *testing.T, conn *sql.DB, id uuid.UUID, email, password string) bool {
	sqlCheck := "SELECT id, email, password FROM users WHERE id = ?"

	row := conn.QueryRow(sqlCheck, id)

	var user user.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		t.Errorf("Error getting user: %s", err)
		return false
	}

	if user.ID != id {
		t.Errorf("Expected ID %s, got %s", id, user.ID)
		return false
	}

	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
		return false
	}

	if user.Password != password {
		t.Errorf("Expected Password %s, got %s", password, user.Password)
		return false
	}

	return true
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testUserMatches(t, db, u.ID, u.Email, u.Password) {
			t.Errorf("User does not match")
		}

	})

	t.Run("error", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Create(&u)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		usr, err := repo.Get(u.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if usr.ID != u.ID {
			t.Errorf("Expected ID %s, got %s", u.ID, usr.ID)
		}

		if usr.Email != u.Email {
			t.Errorf("Expected Email %s, got %s", u.Email, usr.Email)
		}

		if usr.Password != u.Password {
			t.Errorf("Expected Password %s, got %s", u.Password, usr.Password)
		}

	})

	t.Run("error", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, err = repo.Get(uuid.New())

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		u.Email = randomEmail()

		err = repo.Update(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testUserMatches(t, db, u.ID, u.Email, u.Password) {
			t.Errorf("User does not match")
		}

	})

	t.Run("error", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		u.Email = longDummyString(1000)

		err = repo.Update(&u)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "Data too long") {
			t.Errorf("Unexpected error: %s", err)
		}

	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		db, err := sql.Open("mysql", "root:juno@tcp(localhost:3306)/user_test?parseTime=true")

		if err != nil {
			t.Errorf("Error connecting to database: %s", err)
		}

		defer db.Close()

		repo := New(db)

		err = mysql.ExecuteMigrations(db)

		if err != nil {
			t.Errorf("Error executing migrations: %s", err)
		}

		email := randomEmail()

		defer db.Exec("DELETE FROM users WHERE email = ?", email)

		u := user.User{
			ID:       uuid.New(),
			Email:    email,
			Password: "password",
		}

		err = repo.Create(&u)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = repo.Delete(u.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, err = repo.Get(u.ID)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

	})
}
