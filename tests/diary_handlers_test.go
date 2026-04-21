package tests
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"myplants-server/internal/database"
	"myplants-server/internal/handlers"
	"myplants-server/internal/models"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.Content{}, &models.User{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	database.DB = db
	return db
}

func TestCreateAndListDiary(t *testing.T) {
	setupInMemoryDB(t)

	// create a test user directly
	user := models.User{Username: "testuser", Password: "p"}
	database.GetDB().Create(&user)

	// prepare request body
	req := handlers.CreateContentRequest{
		Title:   "My first diary",
		Content: "Planting seeds",
		Images:  []string{"https://example.com/a.jpg"},
		Tags:    []string{"family_Fabaceae", "genus_Vicia"},
	}
	b, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/diaries", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")
	// set authenticated user id in context
	c.Set("userId", uint(user.ID))

	handlers.CreateDiary(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d, body=%s", w.Code, w.Body.String())
	}

	// now list diaries
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/api/diaries?limit=10&offset=0", nil)

	handlers.ListDiaries(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for list, got %d", w2.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse list response: %v", err)
	}
	items := resp["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 diary, got %d", len(items))
	}
}
