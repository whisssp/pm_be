package entity

import "gorm.io/gorm"

//
//import (
//	"database/sql/driver"
//	"errors"
//	"fmt"
//)
//
//type UserRole int
//
//const (
//	RoleAdmin UserRole = iota
//	RoleUser
//)
//
//var AllRoles = [2]string{"ADMIN", "USER"}
//
//func (item *UserRole) String() string {
//	return AllRoles[*item-1]
//}
//
//func parseStr2ItemStatus(s string) (UserRole, error) {
//	for i := range AllRoles {
//		if AllRoles[i] == s {
//			return UserRole(i), nil
//		}
//	}
//
//	return UserRole(0), errors.New("invalid status string")
//}
//
//func (item *UserRole) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprintf("fail to Scan data from database: %s", value))
//	}
//
//	v, err := parseStr2ItemStatus(string(bytes))
//
//	if err != nil {
//		return errors.New(fmt.Sprintf("fail to Scan data from database: %s", value))
//	}
//
//	*item = v
//
//	return nil
//}
//
//func (item *UserRole) Value() (driver.Value, error) {
//	if item == nil {
//		return nil, nil
//	}
//
//	return item.String(), nil
//}

const (
	RoleAdmin int64 = iota + 1
	RoleUser  int64 = iota + 1
)

type UserRole struct {
	gorm.Model
	ID   int64  `gorm:"type:bigint;autoIncrement;primaryKey"`
	Name string `gorm:"type:varchar(50);unique;not null"`
	User []User `gorm:"foreignKey:RoleID"`
}