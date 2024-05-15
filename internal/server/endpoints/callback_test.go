package endpoints

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleCallback(t *testing.T) {
	err := godotenv.Load(".env.test")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " dbname=" + name + " password=" + password + " port=" + port

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.DB()


	db.Migrator().DropTable(&model.GhUser{})
	db.AutoMigrate(&model.GhUser{})

	h := &handlers{db: db}

	request, err := http.NewRequest("GET", "/callback?code=mock-code", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	h.HandleCallback(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var savedUser model.GhUser
	db.First(&savedUser, "username = ?", "szmulinho")
	assert.Equal(t, "doctor", savedUser.Role)
}
