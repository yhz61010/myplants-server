package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "myplants-server/internal/handlers"
)

func TestRegisterAndLogin(t *testing.T) {
    setupInMemoryDB(t)

    // register
    reg := map[string]string{"username": "u1", "password": "passw"}
    b, _ := json.Marshal(reg)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(b))
    c.Request.Header.Set("Content-Type", "application/json")
    handlers.Register(c)
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201 created for register, got %d body=%s", w.Code, w.Body.String())
    }

    // login
    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Request, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
    c2.Request.Header.Set("Content-Type", "application/json")
    handlers.Login(c2)
    if w2.Code != http.StatusOK {
        t.Fatalf("expected 200 OK for login, got %d body=%s", w2.Code, w2.Body.String())
    }
    var resp map[string]interface{}
    if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to parse login response: %v", err)
    }
    if resp["token"] == nil {
        t.Fatalf("expected token in login response")
    }
}
