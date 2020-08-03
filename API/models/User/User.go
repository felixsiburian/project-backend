package User

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"strings"
	"time"
)

type User struct {
	Id         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email      string    `gorm:"size:100; no null; unique" json:"email"`
	Password   string    `gorm:"size:100; not null"json:"password"`
	Role       UserRole  `json:"role"`
	UserRoleId uint32    `json:"user_role_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		log.Fatal(err)
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Id = 0
	u.Email = html.EscapeString(strings.TrimSpace((u.Email)))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		log.Fatal(err)
		return &User{}, err
	}
	if u.UserRoleId != 0 {
		err = db.Debug().Model(&UserRole{}).Where("role_id = ? ", u.UserRoleId).Take(&u.Role).Error
		if err != nil {
			log.Fatal("err ", err)
			return &User{},err
		}
	}
	return u, nil
}

func (u *User) FindAllUser(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		log.Fatal(err)
		return &[]User{}, err
	}

	if len(users) > 0 {
		for i, _ := range users {
			err := db.Debug().Model(&UserRole{}).Where("role_id = ? ", users[i].UserRoleId).Take(&users[i].Role).Error
			if err != nil {
				return &[]User{}, err
			}
		}
	}

	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		log.Fatal(err)
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		log.Fatal(err)
		return &User{}, errors.New("User Not Found")
	}
	if u.UserRoleId != 0 {
		err = db.Debug().Model(&UserRole{}).Where("role_id = ?", u.UserRoleId).Take(&u.Role).Error
		if err != nil {
			return &User{}, err
		}
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (string, error) {
	//hash password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Updates(
		map[string]interface{}{
			"password":  u.Password,
			"email":     u.Email,
			"user_role_id" : u.UserRoleId,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return "Error", db.Error
	}

	//tampilkan update user
	//err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	//if err != nil {
	//	return &User{}, err
	//}
	//
	//if u.UserRoleId != 0 {
	//	err = db.Debug().Model(&UserRole{}).Where("role_id = ?", u.UserRoleId).Take(&u.Role).Error
	//	if err != nil {
	//		return &User{}, err
	//	}
	//}

	return "Update Succesfull", nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		log.Fatal(db.Error)
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
