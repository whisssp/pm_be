package entity

//type Role string
//
//const (
//	RoleAdmin     Role = "admin"
//	RoleModerator Role = "moderator"
//	RoleUser      Role = "user"
//)
//
//func (st *Role) Scan(value interface{}) error {
//	*st = Role(value.([]byte))
//	return nil
//}
//
//func (st *Role) Value() (driver.Value, error) {
//	return string(st), nil
//}
//
//type User struct {
//	gorm.Model
//	Phone    string `gorm:"type:varchar(11)"`
//	Password string `gorm:"type:varchar(255)"`
//	Role     Role   `gorm:"type:role;default:'user'"`
//}