package main

import (
	"api-kasirapp/auth"
	"api-kasirapp/handler"
	"api-kasirapp/helper"
	"api-kasirapp/repository"
	"api-kasirapp/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, username, password, databaseName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := repository.NewRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	productRepository := repository.NewProductRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	supplierRepository := repository.NewSupplierRepository(db)
	discountRepository := repository.NewDiscountRepository(db)

	userService := service.NewService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	productService := service.NewProductService(productRepository)
	customersService := service.NewCustomerService(customerRepository)
	supplierService := service.NewSupplierService(supplierRepository)
	discountService := service.NewDiscountService(discountRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customersService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	discountHandler := handler.NewDiscountHandler(discountService)
	router := gin.Default()

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/categories", categoryHandler.CreateCategory)
	api.POST("/products", productHandler.CreateProduct)
	api.POST("/customers", customerHandler.CreateCustomer)
	api.POST("/suppliers", supplierHandler.CreateSupplier)
	api.POST("/discounts", discountHandler.CreateDiscount)

	api.GET("/categories", categoryHandler.GetCategories)
	api.GET("/categories/:id", categoryHandler.GetCategoryById)
	api.GET("/products", productHandler.GetProducts)
	api.GET("/products/:id", productHandler.GetProductById)
	api.GET("/customers", customerHandler.GetCustomers)
	api.GET("/customers/:id", customerHandler.GetCustomerById)
	api.GET("/suppliers", supplierHandler.GetSuppliers)
	api.GET("/suppliers/:id", supplierHandler.GetSupplierById)
	api.GET("/discounts", discountHandler.GetDiscounts)
	api.GET("/discounts/:id", discountHandler.GetDiscountById)
	api.GET("/category-products/:id", categoryHandler.GetCategoryProducts)

	api.PUT("/categories/:id", categoryHandler.UpdateCategory)
	api.PUT("/products/:id", productHandler.UpdateProduct)
	api.PUT("/customers/:id", customerHandler.UpdateCustomer)
	api.PUT("/suppliers/:id", supplierHandler.UpdateSupplier)
	api.PUT("/discounts/:id", discountHandler.UpdateDiscount)

	api.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	api.DELETE("/products/:id", productHandler.DeleteProduct)
	api.DELETE("/customers/:id", customerHandler.DeleteCustomer)
	api.DELETE("/suppliers/:id", supplierHandler.DeleteSupplier)
	api.DELETE("/discounts/:id", discountHandler.DeleteDiscount)

	router.Run()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func authMiddleware(authService auth.Service, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
