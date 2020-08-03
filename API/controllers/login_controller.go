package controllers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"project-backend/API/auth"
	"project-backend/API/models/User"
	"project-backend/API/responses"
	"project-backend/API/utils/formaterror"
)

func (server *Server) Login (w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := User.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error){
	var err error

	user := User.User{}

	err = server.DB.Debug().Model(User.User{}).Where("email = ?",email).Take(&user).Error
	if err != nil {
		return "Email Not Found", http.ErrBodyNotAllowed
	}

	err = User.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
		return "Password wrong. Try Again", http.ErrBodyNotAllowed
	}

	return auth.CreateToken(user.Id)
}
