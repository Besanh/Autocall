package model

import (
	"time"
)

type Contact struct {
	ID         int       `json:"id" gorm:":id;primaryKey"`
	FullName   string    `json:"fullname" gorm:":fullname;type:varchar(100);not null"`
	Address    string    `json:"address" gorm:":address;type:varchar(255);not null"`
	Phone      string    `json:"phone" gorm:":phone;type:varchar(10);notnull"`
	Status     string    `json:"status" gorm:":status;type:varchar(10);notnull"`
	CallStatus string    `json:"call_status" gorm:":call_status;type:varchar(255);null"`
	Note       string    `json:"note" gorm:":note;type:text;null"`
	CreatedAt  time.Time `json:"created_at" gorm:":created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt  time.Time `json:"updated_at" gorm:":updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	// DeletedAt  time.Time `gorm:",soft_delete"`
	CreatedBy int `json:"created_by" gorm:":created_by,type:integer,notnull,nullzero,default:'0'"`
	UpdatedBy int `json:"updated_by" gorm:":updated_by,type:integer,notnull,nullzero,default:'0'"`
	GroupID   int `json:"group_id" gorm:":group_id,type:integer,notnull,nullzero,default:'0'"`
}

func (Contact) TableName() string {
	return "contact_call_form"
}
