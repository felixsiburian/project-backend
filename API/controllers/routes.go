package controllers

import (
	"project-backend/API/middleware"
)

func (s *Server) InitializeRoutes() {
	//Home ROute
	s.Router.HandleFunc("/", middleware.SetMiddlewareAuthentication(s.Home)).Methods("GET")

	//Login ROute
	s.Router.HandleFunc("/api/login", middleware.SetMiddlewareJSON(s.Login)).Methods("POST")

	// User Route
	s.Router.HandleFunc("/api/users", middleware.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/api/users", middleware.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middleware.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/api/users/{id}", middleware.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")
}