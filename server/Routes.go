package server

import (
	"audioPhile/handlers"
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
		//api.Route("/auth", func(auth chi.Router) {
		//	auth.Use(middlewares.JWTAuthMiddleware)
		//	auth.Route("/users", func(users chi.Router) {
		//
		//	})
		//})
	})

	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
