package handlers

import (
	"kenshi/database"
	"kenshi/models"
	"kenshi/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user.Password = string(hashedPassword)

	// Insert user into database
	result, err := database.GetDB().Collection("user").InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

func SignIn(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	err := database.GetDB().Collection("user").FindOne(c.Context(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := utils.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	user.RefreshToken = refreshToken
	if _, err := database.GetDB().Collection("user").UpdateByID(c.Context(), user.ID, bson.M{"$set": bson.M{"refreshToken": refreshToken}}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not update refresh token"})
	}

	return c.JSON(fiber.Map{
		"token": accessToken,
	})
}

func SignOut(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	_, err := database.GetDB().Collection("user").UpdateByID(c.Context(), userID, bson.M{"$set": bson.M{"refreshToken": ""}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not sign out"})
	}
	return c.SendStatus(fiber.StatusOK)
}
