package model

import (
	"Go_blog/utils/errmsg"
	"gorm.io/gorm"
)

type Category struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// 查询分类是否存在
func CheckCategory(name string) (code int) {
	var cate Category
	db.Select("id").Where("name=?",name).First(&cate)
	if cate.ID>0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCESS
}

// 新增分类
func CreateCategory(data *Category) int {
	err:=db.Create(&data).Error
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询分类信息
func GetCateInfo(id int) (Category,int) {
	var cate Category
	db.Where("id=?",id).First(&cate)
	return cate,errmsg.SUCCESS
}

// 查询分类列表
func GetCategory(pageSize int,pageNum int) ([]Category,int64) {
	var cate []Category
	var total int64
	err = db.Find(&cate).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	db.Model(&cate).Count(&total)
	if err != nil && err!=gorm.ErrRecordNotFound{
		return nil,0
	}
	return cate,total
}

// 编辑分类
func EditCategory(id int,data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"]=data.Name

	err=db.Model(&cate).Where("id=?",id).Updates(maps).Error
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除分类
func DeleteCategory(id int) int {
	var cate Category
	err=db.Where("id=?",id).Delete(&cate).Error
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}