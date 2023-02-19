package orm

import (
	"gorm.io/gorm"
)

type Tbl_member struct {
	gorm.Model
	Id             int
	Ref_code       string
	Child_ref_code string
	Firstname      string
	Lastname       string
	Email          string
	Password       string
}
