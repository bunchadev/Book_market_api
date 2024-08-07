package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("SECRET_KEY")

func GenerateToken(userId string, role string, myTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userId,
		"role":       role,
		"created_at": time.Now(),
		"exp":        time.Now().Add(myTime).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}
func VerifyToken_v1(token string, requiredRoles []string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", errors.New("could not parse token")
	}
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return "", errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("invalid role")
	}
	if check := Contains(role, requiredRoles); !check {
		return "", errors.New("invalid role")
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("invalid userId")
	}
	return userId, nil
}

func Contains(role string, requiredRoles []string) bool {
	for _, r := range requiredRoles {
		if role == r {
			return true
		}
	}
	return false
}

func VerifyToken_v2(token string) (string, time.Duration, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", 0, errors.New("could not parse token")
	}
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return "", 0, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, errors.New("invalid token claims")
	}
	id, ok := claims["userId"].(string)
	if !ok {
		return "", 0, errors.New("invalid userId")
	}
	exp := int64(claims["exp"].(float64))
	expirationTime := time.Unix(exp, 0)
	remainingTime := time.Until(expirationTime)
	return id, remainingTime, nil
}
