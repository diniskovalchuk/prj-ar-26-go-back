package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/diniskovalchuk/prj2/internal/app"
	"github.com/diniskovalchuk/prj2/internal/domain"
	"github.com/diniskovalchuk/prj2/internal/infra/http/controllers"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/upper/db/v4"
)

func AuthMiddleware(ja *jwtauth.JWTAuth, as app.AuthService, us app.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := jwtauth.VerifyRequest(ja, r, jwtauth.TokenFromHeader)

			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				controllers.Unauthorized(w, err)
				return
			}

			var userIDRaw interface{}

			err = token.Get("user_id", &userIDRaw)
			if err != nil {
				controllers.Unauthorized(w, errors.New("missing user_id"))
				return
			}

			var uuidRaw interface{}

			err = token.Get("uuid", &uuidRaw)
			if err != nil {
				controllers.Unauthorized(w, errors.New("missing uuid"))
				return
			}

			userIDFloat, ok := userIDRaw.(float64)
			if !ok {
				controllers.Unauthorized(w, errors.New("invalid user_id"))
				return
			}

			uuidStr, ok := uuidRaw.(string)
			if !ok {
				controllers.Unauthorized(w, errors.New("invalid uuid"))
				return
			}

			uId := uint64(userIDFloat)

			uUuid, err := uuid.Parse(uuidStr)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			auth := domain.Session{
				UserId: uId,
				UUID:   uUuid,
			}
			err = as.Check(auth)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			user, err := us.FindById(uId)
			if err != nil {
				if errors.Is(err, db.ErrNoMoreRows) {
					err = errors.New("unauthorized")
				}
				controllers.Unauthorized(w, err)
				return
			}

			ctx = context.WithValue(ctx, controllers.UserKey, user)
			ctx = context.WithValue(ctx, controllers.SessKey, auth)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
