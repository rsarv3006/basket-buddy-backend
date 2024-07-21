package handler

import (
	"basket-buddy-backend/dto"
	"basket-buddy-backend/ent"
	"basket-buddy-backend/ent/share"
	"basket-buddy-backend/helpers"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateShareEndpoint(dbClient *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*ent.AppUser)

		shareCreateDto := new(dto.CreateShareDto)

		if err := c.BodyParser(shareCreateDto); err != nil {
			return sendBadRequestResponse(c, err, "Failed to parse share data")
		}

		if shareCreateDto.Data == nil {
			return sendBadRequestResponse(c, nil, "share data not defined")
		}

		shareCode, err := createShareCode(dbClient)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		shareObj, err := dbClient.Share.Create().
			SetData(shareCreateDto.Data).
			SetCreatorID(currentUser.ID).
			SetShareCode(shareCode).
			Save(context.Background())

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Share created successfully",
			"shareCode": shareObj.ShareCode,
		})
	}
}

func createShareCode(dbClent *ent.Client) (string, error) {
	shareCode := helpers.GenerateShareCode()

	foundShare, err := dbClent.Share.Query().Where(share.ShareCode(shareCode)).Exist(context.Background())

	if err != nil {
		return "", err
	}

	if foundShare {
		return "", errors.New("Share code already exists")
	}

	return shareCode, nil
}

func FetchShareEndpoint(dbClient *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		shareCode := c.Params("ShareCode")

		if shareCode == "" {
			return sendBadRequestResponse(c, nil, "Share code not defined")
		}

		shareObj, err := dbClient.Share.Query().Where(share.ShareCode(shareCode)).First(context.Background())

		if err != nil {
			if ent.IsNotFound(err) {
				return sendNotFoundResponse(c, err)
			}
			return sendInternalServerErrorResponse(c, err)
		}

		if shareObj.Expiration.Before(time.Now()) {
			// TODO: convert to soft delete and update status to expired
			dbClient.Share.DeleteOne(shareObj).Exec(context.Background())
			return sendNotFoundResponse(c, nil)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":   "Share fetched successfully",
			"shareCode": shareObj.ShareCode,
			"data":      shareObj.Data,
		})
	}
}
