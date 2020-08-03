package User

type UserRole struct {
	RoleId       uint32 `gorm:"primary_key;auto_increment" json:"role_id"`
	UserRole string `json:"user_role"`
}
