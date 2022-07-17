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
		api.Post("/sign_in", handlers.SignInUserHandler)
		api.Post("/sign_up", handlers.SignUpUserHandler)
		api.Route("/auth", func(auth chi.Router) {
			auth.Use(middlewares.JWTAuthMiddleware)
			auth.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.AccessControlMiddleware)
				admin.Post("/sign_up", handlers.AdminSignUpUserHandler)
				admin.Post("/add_product", handlers.AddProductHandler)
				admin.Post("/add_role", handlers.AddAdminRoleHandler)
				admin.Get("/all_users", handlers.FetchAllUsersHandler)
				admin.Post("/upload_asset", handlers.FireBaseUploadHandler)
			})
			auth.Route("/users", func(users chi.Router) {
				users.Post("/add_address", handlers.AddAddressHandler)
				users.Post("/add_card", handlers.AddCardHandler)
				users.Delete("/remove_address", handlers.RemoveAddressHandler)
				users.Delete("/remove_card", handlers.RemoveCardHandler)
				users.Get("/fetch_user", handlers.FetchUserHandler)
			})
			auth.Route("/products", func(products chi.Router) {
				products.Get("/all_products", handlers.FetchProductsHandler)
				products.Get("/fetch_product_category", handlers.FetchProductsCategoryHandler)
				products.Get("/fetch_latest_product", handlers.FetchLatestProductHandler)
				products.Get("/fetch_you_may_like", handlers.FetchYouMayLike)
				products.Get("/fetch_product_assets", handlers.FetchProductAssetsHandler)

			})
			auth.Route("/cart", func(cart chi.Router) {
				cart.Post("/add_cart", handlers.AddCartHandler)
				cart.Post("/add_order", handlers.AddOrderHandler)
				cart.Put("/update_cart_products", handlers.UpdateCartProductsHandler)
				cart.Post("/add_cart_products", handlers.AddCartProductHandler)
			})
		})
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
