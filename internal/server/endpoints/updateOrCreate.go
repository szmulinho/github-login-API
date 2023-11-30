package endpoints

import (
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/gorm"
)

func (h *handlers) updateOrCreateGitHubUser(db *gorm.DB, githubUser model.GithubUser) error {
	existingUser := model.GithubUser{}
	if err := db.Where("login = ?", githubUser.Login).First(&existingUser).Error; err == nil {
		existingUser.Email = githubUser.Email
		return db.Save(&existingUser).Error
	}
	return db.Create(&githubUser).Error
}

func (h *handlers) updateOrCreatePublicRepo(db *gorm.DB, publicRepo model.PublicRepo) error {
	existingRepo := model.PublicRepo{}
	if err := db.Where("name = ?", publicRepo.Name).First(&existingRepo).Error; err == nil {
		existingRepo.Description = publicRepo.Description
		return db.Save(&existingRepo).Error
	}
	return db.Create(&publicRepo).Error
}