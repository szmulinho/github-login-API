package endpoints

import (
	"github.com/google/go-github/github"
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/gorm"
)

func (h *handlers) SaveUserToDatabase(user *github.User, role string, db gorm.DB) error {
	var existingUser model.GithubUser
	db.Where("github_user_id = ?", user.GetID()).First(&existingUser)

	if existingUser.GithubUserID == 0 {
		newUser := model.GithubUser{
			GithubUserID: user.GetID(),
			Name:         user.GetName(),
			Email:        user.GetEmail(),
			Role:         role,
		}
		db.Create(&newUser)
	} else {
		existingUser.Name = user.GetName()
		existingUser.Email = user.GetEmail()
		existingUser.Role = role
		db.Save(&existingUser)
	}

	return nil
}
