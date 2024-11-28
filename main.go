package main

import (
	"api-kasirapp/auth"
	"api-kasirapp/handler"
	"api-kasirapp/helper"
	"api-kasirapp/repository"
	"api-kasirapp/service"
	"fmt"
	"github.com/gin-contrib/cors"
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
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/categories", categoryHandler.CreateCategory)
	api.POST("/products", productHandler.CreateProduct)
	api.POST("/customers", customerHandler.CreateCustomer)
	api.POST("/suppliers", supplierHandler.CreateSupplier)
	api.POST("/discounts", discountHandler.CreateDiscount)
	api.POST("/transactions", transactionHandler.CreateTransaction)
	api.POST("/product-image/:id", productHandler.UploadProductImage)

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
	api.GET("/category-name/:category_name", categoryHandler.GetProductsByCategoryName)

	api.PUT("/categories/:id", categoryHandler.UpdateCategory)
	api.PUT("/products/:id", productHandler.UpdateProduct)
	api.PUT("/customers/:id", customerHandler.UpdateCustomer)
	api.PUT("/suppliers/:id", supplierHandler.UpdateSupplier)
	api.PUT("/discounts/:id", discountHandler.UpdateDiscount)
	api.PUT("/stocks/:id", stockHandler.UpdateStock)

	api.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	api.DELETE("/products/:id", productHandler.DeleteProduct)
	api.DELETE("/customers/:id", customerHandler.DeleteCustomer)
	api.DELETE("/suppliers/:id", supplierHandler.DeleteSupplier)
	api.DELETE("/discounts/:id", discountHandler.DeleteDiscount)
	api.DELETE("/stocks/:id", stockHandler.DeleteStock)

	api.POST("/stocks", stockHandler.AddStock)
	api.GET("/stocks/:id", stockHandler.GetStocks)
	api.GET("/stock-product/:productID", stockHandler.GetStocksByProductID)

	api.GET("/export/products", productHandler.ExportProducts)
	api.POST("/import/products", productHandler.ImportProducts)

	api.GET("/export/customers", customerHandler.ExportCustomers)
	api.POST("/import/customers", customerHandler.ImportCustomers)

	api.GET("/export/suppliers", supplierHandler.ExportSuppliers)
	api.POST("/import/suppliers", supplierHandler.ImportSuppliers)

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
