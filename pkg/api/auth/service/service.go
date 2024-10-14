package service

import (
	"juno/pkg/api/auth"
	"juno/pkg/api/user"
	"juno/pkg/util"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger      logrus.FieldLogger
	userService user.Service
}

func New(logger *logrus.Logger, userService user.Service) *Service {
	return &Service{
		logger:      logger,
		userService: userService,
	}
}

func TokenToUser(token string) (*user.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "Token is expired") {
			return nil, auth.ErrExpiredToken
		}
		return nil, err
	}

	if !t.Valid {
		return nil, auth.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !ok {
		return nil, auth.ErrInvalidToken
	}

	exp, ok := claims["exp"].(float64)

	if !ok {
		return nil, auth.ErrInvalidToken
	}

	if int64(exp) < time.Now().Unix() {
		return nil, auth.ErrExpiredToken
	}

	parseID, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	u := &user.User{
		ID:    parseID,
		Email: claims["email"].(string),
	}

	return u, nil
}

func Token(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    u.ID.String(),
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) Authenticate(email, password string) (string, error) {
	u, err := s.userService.FirstWhereEmail(email)

	if err != nil {
		return "", err
	}

	if u == nil {
		return "", auth.ErrInvalidEmailOrPassword
	}

	if err := util.CompareBcryptPassword(u.Password, password); err != nil {
		return "", auth.ErrInvalidEmailOrPassword
	}

	return Token(u)
}
