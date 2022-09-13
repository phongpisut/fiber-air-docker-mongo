package controllers

import (
	"context"
	"fiber-api/configs"
	"fiber-api/models"
	"fiber-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "register_system")
var settingCollection *mongo.Collection = configs.GetCollection(configs.DB, "settings")

var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.User{
		Id:      primitive.NewObjectID(),
		Name:    user.Name,
		Surname: user.Surname,
		Tel:     user.Tel,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func GetAllUsers(c *fiber.Ctx) error {

	result, err := GetUserAndSetting(c)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"users": result.User, "settings": result.Settings}},
	)

}

func GetUserAndSetting(c *fiber.Ctx) (models.UserAndProperties, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	var settings []models.Properties
	var userAndSetting models.UserAndProperties

	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	resultsSetting, errSetting := settingCollection.Find(ctx, bson.M{})

	if err != nil {
		return userAndSetting, err
	}
	if errSetting != nil {
		return userAndSetting, err
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return userAndSetting, err
		}

		users = append(users, singleUser)
	}

	defer resultsSetting.Close(ctx)
	for resultsSetting.Next(ctx) {
		var property models.Properties
		if err = resultsSetting.Decode(&property); err != nil {
			return userAndSetting, err
		}

		settings = append(settings, property)
	}

	userAndSetting.Settings = settings[0]
	userAndSetting.User = users

	return userAndSetting, nil
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetSetting(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var settings []models.Properties
	defer cancel()

	results, err := settingCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var property models.Properties
		if err = results.Decode(&property); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		settings = append(settings, property)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": settings}},
	)

}

func EditSetting(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var settings models.Properties
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	//validate the request body
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&settings); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"maximum": settings.Maximum}

	result, err := settingCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedMaximum models.Properties
	if result.MatchedCount == 1 {
		err := settingCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedMaximum)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedMaximum}})
}
