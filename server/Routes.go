package server

import (
	"audioPhile/handlers"
	"audioPhile/middlewares"
	"github.com/go-chi/chi"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Post("/sign-in", handlers.SignInUserHandler)
		api.Post("/sign-up", handlers.SignUpUserHandler)
		api.Route("/auth", func(auth chi.Router) {
			auth.Use(middlewares.JWTAuthMiddleware)
			auth.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.AccessControlMiddleware)
				admin.Post("/sign-up", handlers.AdminSignUpUserHandler)
				admin.Post("/add-product", handlers.AddProductHandler)
				admin.Post("/add-role", handlers.AddAdminRoleHandler)
				admin.Get("/all-users", handlers.FetchAllUsersHandler)
				admin.Post("/upload-asset", handlers.FireBaseUploadHandler)
				admin.Post("/add-category", handlers.AddCategoryHandler)
				admin.Get("/fetch-categories", handlers.FetchCategoryHandler)
				admin.Get("/fetch-product-assets", handlers.FetchProductAssetsHandler)
				admin.Get("/fetch-search-products-list", handlers.FetchProductsListSearchHandler)
			})
			auth.Route("/users", func(users chi.Router) {
				users.Post("/add-address", handlers.AddAddressHandler)
				users.Post("/add-card", handlers.AddCardHandler)
				users.Delete("/remove-address", handlers.RemoveAddressHandler)
				users.Delete("/remove-card", handlers.RemoveCardHandler)
				users.Get("/fetch-user", handlers.FetchUserHandler)
			})
			auth.Route("/products", func(products chi.Router) {
				//products.Get("/all-products", handlers.FetchProductsHandler)
				//products.Get("/fetch-product-category", handlers.FetchProductsCategoryHandler)
				//products.Get("/fetch-latest-product", handlers.FetchLatestProductHandler)
				//products.Get("/fetch-you-may-like", handlers.FetchYouMayLike)
				//products.Get("/fetch-categories", handlers.FetchCategoryHandler)
				products.Get("/fetch-product-assets", handlers.FetchProductAssetsHandler)
				products.Get("/fetch-search-products-list", handlers.FetchProductsListSearchHandler)
			})
			auth.Route("/cart", func(cart chi.Router) {
				cart.Post("/add-cart", handlers.AddCartHandler)
				cart.Post("/add-order", handlers.AddOrderHandler)
				cart.Put("/update-cart-products", handlers.UpdateCartProductsHandler)
				cart.Post("/add-cart-products", handlers.AddCartProductHandler)
			})
		})
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
