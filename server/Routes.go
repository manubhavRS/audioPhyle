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
				admin.Post("/add-product", handlers.AddProductHandler)
				admin.Post("/add-role", handlers.FetchAllUsersHandler)
			})
			auth.Route("/users", func(users chi.Router) {
				users.Post("/add-address", handlers.AddAddressHandler)
				users.Post("/add-card", handlers.AddCardHandler)
			})
			auth.Route("/products", func(products chi.Router) {
				products.Post("/all-products", handlers.FetchProductsHandler)
				products.Post("/fetch-product-category", handlers.FetchProductsCategoryHandler)
				products.Post("/fetch-latest-product", handlers.FetchLatestProductHandler)
				products.Post("/fetch-you-may-like", handlers.FetchYouMayLike)

			})
			auth.Route("/cart", func(cart chi.Router) {
				cart.Post("/add-cart", handlers.AddCartHandler)
				cart.Post("/add-order", handlers.AddOrder)
				cart.Post("/update-cart-products", handlers.UpdateCartProducts)
				cart.Post("/add-cart-products", handlers.AddCartProductHandler)
			})
		})
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
