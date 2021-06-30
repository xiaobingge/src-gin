package model

import "dbger/app/common/db"

//import "github.com/jinzhu/gorm"
type Admins struct {
	Id         uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	AppIds     string `gorm:"column:app_ids;not null"`
	Name       string `gorm:"column:name;not null"`
	Avatar     string `gorm:"column:avatar;not null"`
	Email      string `gorm:"column:email;not null"`
	Status     uint64 `gorm:"column:status;default:1;not null"`
}

func ListAdmins(limit,page int) ([]*Admins, uint64,error)  {
	list := make([]*Admins,0)
	var count uint64
	where := "status = 1"
	if err := db.DB.Local.Model(&Admins{}).Where(where).Count(&count).Error; err != nil {
		return list, count ,err
	}

	offset := (page -1) * limit

	if err := db.DB.Local.Where(where).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return  list,count,err
	}
	return list , count ,nil
}