package Employee

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"log"
	"strings"
	"time"
	"unicode"
)

type Employee struct {
	Id          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	NPWP        string    `json:"npwp"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (e *Employee) Prepare() {
	e.Id = 0
	e.Name = html.EscapeString(strings.TrimSpace(e.Name))
	e.Address = html.EscapeString(strings.TrimSpace(e.Address))
	e.PhoneNumber = html.EscapeString(strings.TrimSpace(e.Address))
	e.NPWP = html.EscapeString(strings.TrimSpace(e.NPWP))
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func (e *Employee) Validate() error {
	if e.Name == "" {
		return errors.New("Required Title")
	}
	if e.Address == "" {
		return errors.New("Required Address")
	}
	if len(e.PhoneNumber) < 11 && len(e.PhoneNumber) > 13 {
		return errors.New("You entered Wrong Phone Number")
	}
	//if IsInt(e.PhoneNumber) == false {
	//	return errors.New("Your input is an alphabetic. Try again")
	//}
	//if IsInt(e.NPWP) == false {
	//	return errors.New("Your input is an alphabetic. Try again")
	//}
	return nil
}

func (e *Employee) SaveEmployee(db *gorm.DB) (*Employee, error) {
	var err error
	err = db.Debug().Model(&Employee{}).Create(&e).Error
	if err != nil {
		log.Fatal("error Save Employee", err.Error())
		return &Employee{}, err
	}
	return e, nil
}

func (e *Employee) FindAllEmployee(db *gorm.DB) (*[]Employee, error){
	var err error
	employees:= []Employee{}
	err = db.Debug().Model(&Employee{}).Limit(100).Find(&employees).Error
	if err != nil {
		log.Fatal("Erorr find employee", err.Error())
		return &[]Employee{} ,err
	}
	return &employees, nil
}

func (e *Employee) FindEmployeeById(db *gorm.DB, eid uint32)(*Employee, error){
	var err error
	err = db.Debug().Model(&Employee{}).Where("id = ?", eid).Take(&e).Error
	if err != nil {
		log.Fatal("Error find Employee by ID", err.Error())
		return &Employee{}, err
	}
	return e, nil
}

func (e *Employee) UpdateAnEmployee(db *gorm.DB, eid uint32) (*Employee, error) {
	var err error

	err = db.Debug().Model(&Employee{}).Where("id = ?", eid).Updates(Employee{
		Name: e.Name,
		Address: e.Address,
		PhoneNumber: e.PhoneNumber,
		NPWP: e.NPWP,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		log.Fatal("error update employee : ", err.Error())
		return &Employee{}, err
	}
	return e, nil
}

func (e *Employee) DeleteAnEmployee(db *gorm.DB, eid uint32) (int64, error){
	db = db.Debug().Model(&Employee{}).Where("id = ?", eid).Take(&Employee{}).Delete(&Employee{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error){
			log.Fatal("erorr delete employee", db.Error)
			return 0, errors.New("Employee Not Found")
		}
		return 0,db.Error
	}
	return db.RowsAffected, nil
}
