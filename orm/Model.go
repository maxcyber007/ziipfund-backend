package orm

import (
	"gorm.io/gorm"
)

type Tbl_user struct {
	gorm.Model
	Id             int
	Ref_code       string
	Child_ref_code string
	Firstname      string
	Lastname       string
	Email          string
	Password       string
	Address        string
	City           string
	Province       string
	Zipcode        string
	Phone          string
	Country        string
	Avatar         string
	Privileges     string
}

type Tbl_deposit struct {
	gorm.Model
	Id         int
	User_id    string
	Money      string
	Status     string
	Created_at string
}

type Tbl_profit struct {
	gorm.Model
	Id         int
	Profit     string
	Created_at string
}

type Tbl_bank struct {
	gorm.Model
	Id         int
	User_id    string
	Bank_file  string
	Created_at string
}

type Tbl_document struct {
	gorm.Model
	Id            int
	User_id       string
	Document_file string
	Created_at    string
}
