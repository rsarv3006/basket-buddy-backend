package database

import (
	"basket-buddy-backend/auth"
	"basket-buddy-backend/config"
	"basket-buddy-backend/ent"
	"basket-buddy-backend/ent/appuser"
	"context"
)

func CreateAnonUser(dbClient *ent.Client) (string, error) {
	doesUserExist, err := dbClient.AppUser.
		Query().
		Where(appuser.Role("anon")).
		Exist(context.Background())

	if err != nil {
		return "", err
	}

	if doesUserExist {
		return "", nil
	}

	anonUser, err := dbClient.AppUser.
		Create().
		SetRole("anon").
		SetUsername("anon").
		SetEmail("anon@basketbuddy.com").
		SetIsActive(true).
		Save(context.Background())

	if err != nil {
		return "", err
	}

	jwtSecretString := config.Config("JWT_SECRET")
	token, err := auth.GenerateJWTFromSecret(anonUser, jwtSecretString, auth.FiftyYears)

	return token, err

}
