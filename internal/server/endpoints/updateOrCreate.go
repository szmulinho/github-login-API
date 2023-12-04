package endpoints

import (
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/gorm"
)

func updateOrCreateGitHubUser(db *gorm.DB, githubUser model.GhUser) error {
	existingUser := model.GhUser{}
	if err := db.Where("login = ?", githubUser.Login).First(&existingUser).Error; err == nil {
		if hasChanges := compareUsers(existingUser, githubUser); hasChanges {
			return db.Save(&githubUser).Error
		}
		return nil
	}

	return db.Create(&githubUser).Error
}

func compareUsers(existingUser, newUser model.GhUser) bool {
	return existingUser.AvatarUrl != newUser.AvatarUrl ||
	 existingUser.Role != newUser.Role ||
		existingUser.Login != newUser.Login
}
