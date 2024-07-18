package auth

import (
	"errors"
	"time"

	"basket-buddy-backend/ent"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	AppUser *ent.AppUser
	jwt.StandardClaims
}

const FiftyYears = 50 * 365 * 24 * time.Hour
const SevenDays = 7 * 24 * time.Hour

var (
	ErrExpired = errors.New("token expired")
	ErrInvalid = errors.New("couldn't parse claims")
)

func GenerateJWT(user *ent.AppUser, ctx *fiber.Ctx, tokenDuration time.Duration) (string, error) {
	jwtSecretString := ctx.Locals("JwtSecret").(string)
	return generateToken(user, jwtSecretString, tokenDuration)
}

func GenerateJWTFromSecret(user *ent.AppUser, jwtSecretString string, tokenDuration time.Duration) (string, error) {
	return generateToken(user, jwtSecretString, tokenDuration)
}

func generateToken(user *ent.AppUser, jwtSecretString string, tokenDuration time.Duration) (string, error) {
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

func ValidateToken(signedToken string, ctx *fiber.Ctx) (*ent.AppUser, error) {
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
