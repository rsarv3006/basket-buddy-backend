package handler

import (
	"basket-buddy-backend/dto"
	"basket-buddy-backend/helper"
	"basket-buddy-backend/model"
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/gofiber/fiber/v2"
)

var BasketBuddyShareCollectionName = "BasketBuddy-Share"

func CreateShareEndpoint(client *firestore.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(map[string]interface{})

		shareCreateDto := new(dto.CreateShareDto)

		if err := c.BodyParser(shareCreateDto); err != nil {
			return sendBadRequestResponse(c, err, "Failed to parse share data")
		}

		if shareCreateDto.Data == nil {
			return sendBadRequestResponse(c, nil, "share data not defined")
		}

		shareCode, err := createShareCode(client)
		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		shareObj := map[string]interface{}{
			"data":       shareCreateDto.Data,
			"expiration": time.Now().Add(time.Hour * 24 * 7),
			"creator_id": currentUser["id"],
			"share_code": shareCode,
			"status":     model.ShareStatusCreated,
			"created_at": firestore.ServerTimestamp,
		}
		docRef, _, err := client.Collection(BasketBuddyShareCollectionName).Add(context.Background(), shareObj)
		if err != nil && err.Error() != "no more items in iterator" {
			log.Println(err)
			return sendInternalServerErrorResponse(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Share created successfully",
			"shareCode": shareCode,
			"id":        docRef.ID,
		})
	}
}

func createShareCode(client *firestore.Client) (string, error) {
	ctx := context.Background()
	shareCode := helper.GenerateShareCode()
	iter := client.Collection(BasketBuddyShareCollectionName).Where("share_code", "==", shareCode).Documents(ctx)
	doc, err := iter.Next()
	if err == nil && doc.Exists() {
		return "", errors.New("Share code already exists")
	}
	if err != nil && err.Error() != "no more items in iterator" {
		return "", err
	}
	return shareCode, nil
}

func FetchShareEndpoint(client *firestore.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		shareCode := c.Params("ShareCode")

		if shareCode == "" {
			return sendBadRequestResponse(c, nil, "Share code not defined")
		}

		ctx := context.Background()
		iter := client.Collection(BasketBuddyShareCollectionName).Where("share_code", "==", shareCode).Where("status", "==", model.ShareStatusCreated).Documents(ctx)
		doc, err := iter.Next()
		if err != nil || !doc.Exists() {
			return sendNotFoundResponse(c, err)
		}

		shareObj := doc.Data()
		expiration, ok := shareObj["expiration"].(time.Time)
		if !ok || expiration.Before(time.Now()) {
			// Mark as expired
			_, err := doc.Ref.Update(ctx, []firestore.Update{{Path: "status", Value: model.ShareStatusExpired}})
			if err != nil {
				return sendInternalServerErrorResponse(c, err)
			}
			return sendNotFoundResponse(c, nil)
		}

		// Mark as accessed
		_, err = doc.Ref.Update(ctx, []firestore.Update{{Path: "status", Value: model.ShareStatusAccessed}})
		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":   "Share fetched successfully",
			"shareCode": shareObj["share_code"],
			"data":      shareObj["data"],
		})
	}
}
