package appmiddleware

import (
	"context"
	"go-restapi/pkg/common"
	"go-restapi/pkg/common/token"
	"net/http"
	"strings"
)

type Contextkey string

const UserContextKey Contextkey = "user"

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.RespondWithError(w, common.ErrTokenMissing)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			common.RespondWithError(w, common.ErrTokenInvalid)
			return
		}

		if strings.ToLower(parts[0]) != "bearer" {
			common.RespondWithError(w, common.ErrTokenInvalid)
			return
		}

		tokenString := parts[1]

		payload, err := token.VerifyToken(tokenString)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
