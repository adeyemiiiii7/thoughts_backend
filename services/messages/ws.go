package messages

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"

	"thoughts_backend_api/models"
	apptypes "thoughts_backend_api/types"
)

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	user, err := h.authenticateWebSocketRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*", "127.0.0.1:*"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	h.hub.Add(user.ID, conn)
	defer h.hub.Remove(user.ID, conn)

	ctx := conn.CloseRead(context.Background())
	_ = wsjson.Write(ctx, conn, map[string]any{
		"type":    "connection.ready",
		"user_id": user.ID,
	})

	<-ctx.Done()
	_ = conn.Close(websocket.StatusNormalClosure, "closing connection")
}

func (h *Handler) authenticateWebSocketRequest(r *http.Request) (*models.User, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	claims := &apptypes.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return h.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
