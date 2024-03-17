package legler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kashyab12/gator/internal/database"
	"net/http"
	"time"
)

type PostUserBody struct {
	Name string `json:"name"`
}

func (config *ApiConfig) PostUsersLegler(w http.ResponseWriter, r *http.Request) {
	var (
		decoder  = json.NewDecoder(r.Body)
		userBody = PostUserBody{}
	)
	if decodeErr := decoder.Decode(&userBody); decodeErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, decodeErr.Error())
	} else if newUser, userInsertErr := config.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userBody.Name,
	}); userInsertErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, userInsertErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusCreated, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}{newUser.ID, newUser.CreatedAt, newUser.UpdatedAt, newUser.Name}); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}
