package utils

import (
	"encoding/json"
	"log"
	"log/slog"
	"mangathorg/internal/models"
	"os"
	"time"
)

var jsonFile = Path + "data/users.json"
var TempUsers []models.TempUser

// retrieveUsers
// retrieves all models.User present in jsonFile and stores them in a slice of models.User.
// It returns the slice of models.User and an error.
func retrieveUsers() ([]models.User, error) {
	var users []models.User

	data, err := os.ReadFile(jsonFile)

	if len(data) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// changeUsers
// overwrites jsonFile with `users` in json format.
func changeUsers(users []models.User) {
	data, errJSON := json.MarshalIndent(users, "", "\t")
	if errJSON != nil {
		Logger.Error(GetCurrentFuncName()+" JSON MarshalIndent error!", slog.Any("output", errJSON))
		return
	}
	errWrite := os.WriteFile(jsonFile, data, 0666)
	if errWrite != nil {
		Logger.Error(GetCurrentFuncName()+" WriteFile error!", slog.Any("output", errWrite))
	}
}

// GetIdNewUser
// returns first unused id in jsonFile.
func GetIdNewUser() int {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	var id int
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, user := range users {
			if user.Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}

// CreateUser
// adds the models.User `newUser` to jsonFile.
func CreateUser(newUser models.User) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	users = append(users, newUser)
	changeUsers(users)
}

// removeUser
// remove the models.User which models.User.Id is sent in argument from jsonFile.
func removeUser(id int) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
	changeUsers(users)
}

// SelectUser
// returns the models.User which models.User.Username matches the `username` argument.
func SelectUser(username string) (models.User, bool) {
	var user models.User
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	var ok bool
	for _, singleUser := range users {
		if singleUser.Username == username {
			ok = true
			user = singleUser
		}
	}
	return user, ok
}

// updateUser
// modifies the models.User in jsonFile that matches
// `updatedUser`'s Id with `updatedUser`'s content.
func updateUser(updatedUser models.User) {
	users, err := retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == updatedUser.Id {
			users[i] = updatedUser
		}
	}
	changeUsers(users)
}

func deleteTempUser(temp models.TempUser) {
	for i, user := range TempUsers {
		if user == temp {
			TempUsers = append(TempUsers[:i], TempUsers[i+1:]...)
		}
	}
}

func PushTempUser(id string) {
	log.Printf("TempUsers: %#v\n", TempUsers)
	log.Printf("id: %#v\n", id)
	for _, temp := range TempUsers {
		if temp.ConfirmID == id {
			temp.User.Id = GetIdNewUser()
			CreateUser(temp.User)
			deleteTempUser(temp)
		}
	}
}

func ManageTempUsers() {
	duration := setDailyTimer()
	for {
		for _, user := range TempUsers {
			if time.Now().Sub(user.CreationTime) > time.Hour*12 {
				deleteTempUser(user)
			}
		}
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}
