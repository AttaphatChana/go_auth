package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtService struct{ Secret []byte }

func (s *JwtService) Generate(userID string) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID, "exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return tok.SignedString(s.Secret)
}
func (s *JwtService) Validate(token string) (string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.Secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", err
	}
	claims := parsed.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
