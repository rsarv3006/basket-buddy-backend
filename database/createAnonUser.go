package database

import (
	"basket-buddy-backend/auth"
	"basket-buddy-backend/config"
	"basket-buddy-backend/ent"
	"basket-buddy-backend/ent/appuser"
	"context"
	"log"
)

func CreateAnonUser(dbClient *ent.Client) {
	doesUserExist, err := dbClient.AppUser.
		Query().
		Where(appuser.Role("anon")).
		Exist(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	if doesUserExist {
		log.Println("Anon user already exists")
		return
	}

	anonUser, err := dbClient.AppUser.
		Create().
		SetRole("anon").
		SetUsername("anon").
		SetEmail("anon@basketbuddy.com").
		SetIsActive(true).
		Save(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	jwtSecretString := config.Config("JWT_SECRET")
	token, err := auth.GenerateJWTFromSecret(anonUser, jwtSecretString, auth.FiftyYears)

	log.Println("Created anon user with token: " + token)

	if err != nil {
		log.Fatal(err)
	}
}
