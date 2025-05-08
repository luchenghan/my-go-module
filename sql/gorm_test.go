package sql

import (
	"fmt"
	"log"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gdb *gorm.DB

func init() {
	var err error
	db, err = NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dial := mysql.New(mysql.Config{
		Conn: db,
	})

	gdb, err = gorm.Open(dial)
	if err != nil {
		log.Fatalf("Failed to open gorm: %v", err)
	}
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
}

func Test_Gorm_Create(t *testing.T) {
	u := new(User)
	u.Username = "test"
	u.Email = "test"

	tx := gdb.Create(u)
	if tx.Error != nil {
		log.Fatalf("Failed to create user: %v", tx.Error)
	}
	log.Println("User created successfully")
}

func Test_Gorm_Find(t *testing.T) {
	u := new(User)
	tx := gdb.Where("email = ?", "test").First(u)
	if tx.Error != nil {
		log.Fatalf("Failed to find user: %v", tx.Error)
	}
	log.Printf("User ID: %d, Username: %s, Email: %s", u.ID, u.Username, u.Email)
}

func Test_Gorm_FindAll(t *testing.T) {
	var users []User
	tx := gdb.Where("username = ?", "test").Find(&users)
	if tx.Error != nil {
		log.Fatalf("Failed to find users: %v", tx.Error)
	}
	for _, user := range users {
		log.Printf("User ID: %d, Username: %s, Email: %s", user.ID, user.Username, user.Email)
	}
}

func Test_Gorm_Update(t *testing.T) {
	u := new(User)
	u.ID = 1020
	u.Username = "test"
	u.Email = "test"
	tx := gdb.Save(u)
	if tx.Error != nil {
		log.Fatalf("Failed to update user: %v", tx.Error)
	}
	log.Println("User updated successfully")
}

func Test_Gorm_FindLastIDAndDelete(t *testing.T) {
	u := new(User)

	tx := gdb.Order("id desc").Limit(1).Find(u)
	if tx.Error != nil {
		log.Fatalf("Failed to find user: %v", tx.Error)
	}

	tx = gdb.Delete(u)
	if tx.Error != nil {
		log.Fatalf("Failed to delete user: %v", tx.Error)
	}
	log.Println("User deleted successfully")
}

func Test_Gorm_Transaction(t *testing.T) {
	tx := gdb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back")
		}
	}()

	u := new(User)
	u.Username = "test"
	u.Email = "test"

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to create user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
	log.Println("Transaction committed successfully")
}

func Test_Gorm_BatchInsert(t *testing.T) {
	start := time.Now()
	defer func() {
		log.Printf("Batch insert took %s", time.Since(start))
	}()

	var users []*User
	for i := 0; i < 1000; i++ {
		users = append(users, &User{
			Username: fmt.Sprintf("gorm_value%d", i),
			Email:    fmt.Sprintf("gorm_value%d", i+1),
		})
	}

	tx := gdb.CreateInBatches(users, 100)
	if tx.Error != nil {
		log.Fatalf("Failed to batch insert users: %v", tx.Error)
	}
	log.Println("Batch insert successful")
}
