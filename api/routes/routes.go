package routes

import (
	"fiber-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {

	app.Post("/addUser", controllers.CreateUser)
	app.Get("/userList", controllers.GetAllUsers)

}

func AdminRoute(app *fiber.App) {
	app.Get("/getSetting", controllers.GetSetting)
	app.Put("/editSetting/:id", controllers.EditSetting)
	app.Delete("/removeUser/:userId", controllers.DeleteAUser)
}
