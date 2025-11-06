package database

import (
	"basket-buddy-backend/auth"
	"basket-buddy-backend/config"
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

var BasketBuddyAppUserCollectionName = "BasketBuddy-AppUser"

func CreateAnonUser(client *firestore.Client) {
	log.Println("Creating Anon User")
	ctx := context.Background()
	// Query for anon user
	log.Println("Query for user")
	iter := client.Collection(BasketBuddyAppUserCollectionName).Where("role", "==", "anon").Documents(ctx)
	doc, err := iter.Next()

	if err == nil && doc.Exists() {
		log.Println("Anon user already exists")
		return
	}

	// Create anon user
	anonUser := map[string]interface{}{
		"role":       "anon",
		"username":   "anon",
		"email":      "anon@basketbuddy.com",
		"is_active":  true,
		"created_at": firestore.ServerTimestamp,
	}
	log.Println("add user")
	docRef, _, err := client.Collection(BasketBuddyAppUserCollectionName).Add(ctx, anonUser)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare user data for JWT
	anonUser["id"] = docRef.ID
	jwtSecretString := config.Config("JWT_SECRET")
	token, err := auth.GenerateJWTFromSecret(anonUser, jwtSecretString, auth.FiftyYears)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created anon user with token: " + token)
}
