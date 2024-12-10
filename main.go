package main

import (
	"api-kasirapp/auth"
	"api-kasirapp/handler"
	"api-kasirapp/helper"
	"api-kasirapp/repository"
	"api-kasirapp/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
	secretKey := os.Getenv("SECRET_KEY")

	// Initialize the auth service
	authService := auth.NewService(secretKey)

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
	stockRepository := repository.NewStockRepository(db)
	transactionRepository := repository.NewOrderRepository(db)

	userService := service.NewService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	productService := service.NewProductService(productRepository, categoryRepository)
	customersService := service.NewCustomerService(customerRepository)
	supplierService := service.NewSupplierService(supplierRepository)
	discountService := service.NewDiscountService(discountRepository)
	stockService := service.NewStockService(stockRepository, productRepository)
	transactionService := service.NewOrderService(transactionRepository, productRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customersService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	discountHandler := handler.NewDiscountHandler(discountService)
	stockHandler := handler.NewStockHandler(stockService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checkers", userHandler.CheckEmailAvailability)
	api.POST("/categories", authMiddleware(authService, userService), categoryHandler.CreateCategory)
	api.POST("/products", authMiddleware(authService, userService), productHandler.CreateProduct)
	api.POST("/customers", authMiddleware(authService, userService), customerHandler.CreateCustomer)
	api.POST("/suppliers", authMiddleware(authService, userService), supplierHandler.CreateSupplier)
	api.POST("/discounts", authMiddleware(authService, userService), discountHandler.CreateDiscount)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/product-image/:id", authMiddleware(authService, userService), productHandler.UploadProductImage)

	api.GET("/categories", authMiddleware(authService, userService), categoryHandler.GetCategories)
	api.GET("/categories/:id", authMiddleware(authService, userService), categoryHandler.GetCategoryById)
	api.GET("/products", authMiddleware(authService, userService), productHandler.GetProducts)
	api.GET("/products/:id", authMiddleware(authService, userService), productHandler.GetProductById)
	api.GET("/customers", authMiddleware(authService, userService), customerHandler.GetCustomers)
	api.GET("/customers/:id", authMiddleware(authService, userService), customerHandler.GetCustomerById)
	api.GET("/suppliers", authMiddleware(authService, userService), supplierHandler.GetSuppliers)
	api.GET("/suppliers/:id", authMiddleware(authService, userService), supplierHandler.GetSupplierById)
	api.GET("/discounts", authMiddleware(authService, userService), discountHandler.GetDiscounts)
	api.GET("/discounts/:id", authMiddleware(authService, userService), discountHandler.GetDiscountById)
	api.GET("/category-products/:id", authMiddleware(authService, userService), categoryHandler.GetCategoryProducts)
	api.GET("/category-name/:category-name", authMiddleware(authService, userService), categoryHandler.GetProductsByCategoryName)

	api.PUT("/categories/:id", authMiddleware(authService, userService), categoryHandler.UpdateCategory)
	api.PUT("/products/:id", authMiddleware(authService, userService), productHandler.UpdateProduct)
	api.PUT("/customers/:id", authMiddleware(authService, userService), customerHandler.UpdateCustomer)
	api.PUT("/suppliers/:id", authMiddleware(authService, userService), supplierHandler.UpdateSupplier)
	api.PUT("/discounts/:id", authMiddleware(authService, userService), discountHandler.UpdateDiscount)
	api.PUT("/stocks/:id", authMiddleware(authService, userService), stockHandler.UpdateStock)

	api.DELETE("/categories/:id", authMiddleware(authService, userService), categoryHandler.DeleteCategory)
	api.DELETE("/products/:id", authMiddleware(authService, userService), productHandler.DeleteProduct)
	api.DELETE("/customers/:id", authMiddleware(authService, userService), customerHandler.DeleteCustomer)
	api.DELETE("/suppliers/:id", authMiddleware(authService, userService), supplierHandler.DeleteSupplier)
	api.DELETE("/discounts/:id", authMiddleware(authService, userService), discountHandler.DeleteDiscount)
	api.DELETE("/stocks/:id", authMiddleware(authService, userService), stockHandler.DeleteStock)

	api.POST("/stocks", authMiddleware(authService, userService), stockHandler.AddStock)
	api.GET("/stocks/:id", authMiddleware(authService, userService), stockHandler.GetStocksByStockID)
	api.GET("/stocks", authMiddleware(authService, userService), stockHandler.GetStocks)
	api.GET("/stock-product/:productID", authMiddleware(authService, userService), stockHandler.GetStocksByProductID)

	api.GET("/export/products", authMiddleware(authService, userService), productHandler.ExportProducts)
	api.POST("/import/products", authMiddleware(authService, userService), productHandler.ImportProducts)

	api.GET("/export/customers", authMiddleware(authService, userService), customerHandler.ExportCustomers)
	api.POST("/import/customers", authMiddleware(authService, userService), customerHandler.ImportCustomers)

	api.GET("/export/suppliers", authMiddleware(authService, userService), supplierHandler.ExportSuppliers)
	api.POST("/import/suppliers", authMiddleware(authService, userService), supplierHandler.ImportSuppliers)

	err = router.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
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
