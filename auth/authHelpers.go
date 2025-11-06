package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	AppUser map[string]interface{}
	jwt.StandardClaims
}

const FiftyYears = 50 * 365 * 24 * time.Hour
const SevenDays = 7 * 24 * time.Hour

var (
	ErrExpired = errors.New("token expired")
	ErrInvalid = errors.New("couldn't parse claims")
)

func GenerateJWT(user map[string]interface{}, ctx *fiber.Ctx, tokenDuration time.Duration) (string, error) {
	jwtSecretString := ctx.Locals("JwtSecret").(string)
	return generateToken(user, jwtSecretString, tokenDuration)
}

func GenerateJWTFromSecret(user map[string]interface{}, jwtSecretString string, tokenDuration time.Duration) (string, error) {
	return generateToken(user, jwtSecretString, tokenDuration)
}

func generateToken(user map[string]interface{}, jwtSecretString string, tokenDuration time.Duration) (string, error) {
	jwtKey := []byte(jwtSecretString)
	expirationTime := time.Now().Add(tokenDuration)
	claims := &JWTClaims{
		AppUser: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateToken(signedToken string, ctx *fiber.Ctx) (map[string]interface{}, error) {
	jwtSecretString := ctx.Locals("JwtSecret").(string)
	jwtKey := []byte(jwtSecretString)
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalid
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrExpired
	}
	return claims.AppUser, nil
}

