package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int64     `json:"id"`
	Uuid      uuid.UUID `json:"uuid"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreateAt  string    `json:"create_at"`
	UpdateAt  string    `json:"update_at"`
}

func main() {
	db, err := sql.Open("postgres", "user=root password=root dbname=app-go sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert multiple users with fake data
	fakeUsers := []User{
		{
			Uuid:      uuid.New(),
			Email:     "user1@example.com",
			Password:  "hashed_password_1",
			FirstName: "Alice",
			LastName:  "Johnson",
		},
		{
			Uuid:      uuid.New(),
			Email:     "user2@example.com",
			Password:  "hashed_password_2",
			FirstName: "Bob",
			LastName:  "Smith",
		},
		{
			Uuid:      uuid.New(),
			Email:     "user3@example.com",
			Password:  "hashed_password_3",
			FirstName: "Charlie",
			LastName:  "Brown",
		},
	}

	for _, user := range fakeUsers {
		_, err := db.Exec(`
			INSERT INTO users (uuid, email, password, first_name, last_name)
			VALUES ($1, $2, $3, $4, $5)`,
			user.Uuid, user.Email, user.Password, user.FirstName, user.LastName)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Multiple users inserted successfully!")
}

func alterTable(db *sql.DB) {
	_, err := db.Exec(`ALTER TABLE users `)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Multiple users inserted successfully!")

}
