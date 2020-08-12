package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"project-backend/API/models/Employee"
	"project-backend/API/responses"
	"project-backend/API/utils/formaterror"
	"strconv"
)

func (server *Server) CreateEmployee(w http.ResponseWriter, r *http.Request){
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee := Employee.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee.Prepare()
	err = employee.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	employeeCreated, err := employee.SaveEmployee(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, employeeCreated.Id))
	responses.JSON(w, http.StatusCreated, employeeCreated)
}

func (server *Server)GetAllEmployee(w http.ResponseWriter, r *http.Request){
	employee := Employee.Employee{}

	employees, err := employee.FindAllEmployee(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, employees)
}

func (server *Server) GetAnEmployee(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	eid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	employee := Employee.Employee{}
	employeeGotten, err := employee.FindEmployeeById(server.DB, uint32(eid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, employeeGotten)
}

func (server *Server) UpdateEmployee(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	eid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee := Employee.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	employee.Prepare()
	err = employee.Validate()
	if err != nil{
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedEmployee, err := employee.UpdateAnEmployee(server.DB, uint32((eid)))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedEmployee)
}

func (server *Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars  := mux.Vars(r)

	employee := Employee.Employee{}

	eid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = employee.DeleteAnEmployee(server.DB, uint32(eid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entitt", fmt.Sprintf("%d", eid))
	responses.JSON(w, http.StatusNoContent, "")
}
