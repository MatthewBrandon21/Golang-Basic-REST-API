package main

import (
	"basic-rest-api-jwt-mysql/config"
	"basic-rest-api-jwt-mysql/controller"
	"basic-rest-api-jwt-mysql/middleware"
	"basic-rest-api-jwt-mysql/repository"
	"basic-rest-api-jwt-mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	itemRepository repository.ItemRepository = repository.NewItemRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	itemService    service.ItemService       = service.NewItemService(itemRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	itemController controller.ItemController = controller.NewItemController(itemService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	authRoutes := r.Group("auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	itemRoutes := r.Group("items", middleware.AuthorizeJWT(jwtService))
	{
		itemRoutes.GET("/", itemController.All)
		itemRoutes.POST("/", itemController.Insert)
		itemRoutes.GET("/:id", itemController.FindByID)
		itemRoutes.PUT("/:id", itemController.Update)
		itemRoutes.DELETE("/:id", itemController.Delete)
	}

	r.Run()
}