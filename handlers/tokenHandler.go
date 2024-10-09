package handlers

import (
	"kenshi/database"
	"kenshi/models"
	"kenshi/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Locals("token")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Missing refresh token"})
	}

	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	// Find the user by userID and validate stored refresh token
	var user models.User

	if err := database.GetDB().Collection("user").FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "User not found"})
	}
	if user.RefreshToken != refreshToken {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid refresh token"})
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{"token": accessToken})
}
