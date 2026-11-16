package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"marketplace/database"
	"marketplace/handlers"
	"marketplace/middleware"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("==================================================")
	log.Println("🛒 Онлайн-Маркетплейс (Multi-Vendor Platform)")
	log.Println("🏢 Microsoft SQL Server & Go REST API")
	log.Println("==================================================")

	dsn := os.Getenv("MARKETPLACE_DSN")
	if dsn == "" {
		dsn = "server=localhost;port=14330;user id=Shahzod;password=1012;database=MarketplaceDB;encrypt=disable"
	}

	database.Connect(dsn)

	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Папка uploads: %v", err)
	}

	r := mux.NewRouter()

	uploadsPath, _ := filepath.Abs(uploadsDir)
	fs := http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsPath)))
	r.PathPrefix("/uploads/").Handler(fs)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	api.HandleFunc("/auth/login", handlers.Login).Methods("POST")

	api.HandleFunc("/products", handlers.ListPublicProducts).Methods("GET")
	api.HandleFunc("/products/{idOrSlug}", handlers.GetPublicProduct).Methods("GET")
	api.HandleFunc("/categories", handlers.ListCategories).Methods("GET")

	api.HandleFunc("/shops", handlers.ListShops).Methods("GET")
	api.HandleFunc("/shops/{idOrSlug}", handlers.GetShop).Methods("GET")

	api.HandleFunc("/reviews/{productId}", handlers.GetProductReviews).Methods("GET")

	authRouter := api.PathPrefix("").Subrouter()
	authRouter.Use(middleware.IsAuthorized)

	authRouter.HandleFunc("/auth/me", handlers.GetMe).Methods("GET")
	authRouter.HandleFunc("/upload", handlers.UploadFile).Methods("POST")

	authRouter.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
	authRouter.HandleFunc("/orders/my", handlers.GetBuyerOrders).Methods("GET")

	authRouter.HandleFunc("/reviews", handlers.CreateReview).Methods("POST")

	sellerRouter := api.PathPrefix("/seller").Subrouter()
	sellerRouter.Use(middleware.IsAuthorized)
	sellerRouter.Use(middleware.RequireSeller)

	sellerRouter.HandleFunc("/products", handlers.SellerListProducts).Methods("GET")
	sellerRouter.HandleFunc("/products", handlers.SellerCreateProduct).Methods("POST")
	sellerRouter.HandleFunc("/products/{id:[0-9]+}", handlers.SellerUpdateProduct).Methods("PUT")
	sellerRouter.HandleFunc("/products/{id:[0-9]+}", handlers.SellerDeleteProduct).Methods("DELETE")

	sellerRouter.HandleFunc("/orders", handlers.GetSellerSubOrders).Methods("GET")
	sellerRouter.HandleFunc("/orders/{id:[0-9]+}/status", handlers.UpdateSubOrderStatus).Methods("PUT")

	sellerRouter.HandleFunc("/stats/monthly", handlers.GetSellerMonthlyStats).Methods("GET")

	sellerRouter.HandleFunc("/my-shop", handlers.GetMyShop).Methods("GET")
	sellerRouter.HandleFunc("/my-shop", handlers.UpdateMyShop).Methods("PUT")

	modRouter := api.PathPrefix("/moderator").Subrouter()
	modRouter.Use(middleware.IsAuthorized)
	modRouter.Use(middleware.RequireModerator)

	modRouter.HandleFunc("/products/pending", handlers.ListPendingProducts).Methods("GET")
	modRouter.HandleFunc("/products", handlers.ListAllProductsModerator).Methods("GET")
	modRouter.HandleFunc("/products/{id:[0-9]+}/review", handlers.ReviewProduct).Methods("PUT")

	modRouter.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")
	modRouter.HandleFunc("/categories/{id:[0-9]+}", handlers.UpdateCategory).Methods("PUT")
	modRouter.HandleFunc("/categories/{id:[0-9]+}", handlers.DeleteCategory).Methods("DELETE")

	modRouter.HandleFunc("/stats", handlers.GetModeratorStats).Methods("GET")

	corsHandler := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"}),
		gorillaHandlers.ExposedHeaders([]string{"Content-Length"}),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8089"
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("🚀 Сервер маркетплейса запущен на http://localhost%s", port)
	log.Printf("📡 API доступно по адресу: http://localhost%s/api", port)
	log.Printf("📁 Раздача файлов: http://localhost%s/uploads/", port)

	if err := http.ListenAndServe(port, corsHandler(r)); err != nil {
		log.Fatalf("❌ Ошибка запуска HTTP сервера: %v", err)
	}
}
