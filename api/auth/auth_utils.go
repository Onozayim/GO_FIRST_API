package auth

import (
	"api/models"
	ctypes "api/types"
	"api/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("error creating token")
	}

	return string(hash), nil
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"id":       user.Id,
		"exp":      time.Now().Add(time.Hour * (24 * 30)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETSTRING")))

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func CheckAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromRequest(r)
		if err != nil {
			utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
			return
		}

		token, err := validateToken(tokenString)
		if err != nil {
			utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
			return
		}

		if !token.Valid {
			utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ReturnErrorStatus(err, http.StatusBadRequest, w)
			return
		}

		// "username": user.Username,
		// "email":    user.Email,
		// "id":       user.Id,
		// "exp":      time.Now().Add(time.Hour * (24 * 30)).Unix(),
		user := models.User{}
		user.Username = claims["username"].(string)
		user.Id = int64(claims["id"].(float64))
		user.Email = claims["email"].(string)

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctypes.UserIDKey, user.Id)
		ctx = context.WithValue(ctx, ctypes.UserUserNameKey, user.Username)
		ctx = context.WithValue(ctx, ctypes.UserEmailKey, user.Email)

		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth, nil
	}
	return "", fmt.Errorf("no se entrego un token")
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("hubo un error inesperado")
		// }
		return []byte(os.Getenv("SECRETSTRING")), nil
	})
}

func GetUserNameFromContext(ctx context.Context, user *models.User) error {
	user_id, ok := ctx.Value(ctypes.UserIDKey).(int64)

	if !ok {
		return fmt.Errorf("ocurrio un error inesperado")
	}
	user_username, ok := ctx.Value(ctypes.UserUserNameKey).(string)

	if !ok {
		return fmt.Errorf("ocurrio un error inesperado")
	}

	user_email, ok := ctx.Value(ctypes.UserEmailKey).(string)

	if !ok {
		return fmt.Errorf("ocurrio un error inesperado")
	}
	user.Id = user_id
	user.Email = user_email
	user.Username = user_username

	return nil
}
