package controllers

import (
	"net/http"
	"project-backend/API/responses"
)

func (server *Server) Home (w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to this Awesome Project")
}