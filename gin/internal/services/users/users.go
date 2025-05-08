package users

import (
	"mymodule/gin/internal/models"
	"mymodule/gin/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var usersRepo *repository.UsersRepository

func Init() {
	usersRepo = repository.NewUsersRepository()
}

func GetUserByID(id string) (*models.User, error) {
	return usersRepo.GetUserByID(id)
}

func GetUsers() ([]models.User, error) {
	return usersRepo.GetAllUsers()
}

func CreateUser(user models.User) error {
	// Hash the password before saving
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return usersRepo.SaveUser(user)
}

func AuthenticateUser(email, password string) (*models.User, error) {
	// Fetch user by email (add this to repository)
	user, err := usersRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare password (assumes password is hashed in DB)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
