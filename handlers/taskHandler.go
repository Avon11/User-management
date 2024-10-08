package handlers

import (
	"time"

	"kenshi/database"
	"kenshi/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	task.UserID = userID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	result, err := database.GetDB().Collection("task").InsertOne(c.Context(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(task)
}

func GetAllTasks(c *fiber.Ctx) error {
	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	cursor, err := database.GetDB().Collection("task").Find(c.Context(), bson.M{"userId": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tasks",
		})
	}
	defer cursor.Close(c.Context())

	var tasks []models.Task
	if err := cursor.All(c.Context(), &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tasks",
		})
	}

	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	var task models.Task
	err = database.GetDB().Collection("task").FindOne(c.Context(), bson.M{"_id": taskID, "userId": userID}).Decode(&task)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UpdatedAt:   time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"status":      updatedTask.Status,
			"updatedAt":   updatedTask.UpdatedAt,
		},
	}

	updatedResp, err := database.GetDB().Collection("task").UpdateByID(c.Context(), bson.M{"_id": taskID, "userId": userID}, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating task",
		})
	}

	if updatedResp.ModifiedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteTask(c *fiber.Ctx) error {
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	userID, _ := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	deleteObj, err := database.GetDB().Collection("task").DeleteOne(c.Context(), bson.M{"_id": taskID, "userId": userID})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting task",
		})
	}

	if deleteObj.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
