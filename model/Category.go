package model

import (
	"Go_blog/utils/errmsg"
	"gorm.io/gorm"
	"errors"
)

type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// CheckCategory 查询分类是否存在
// @name 传过来的name字符串
func CheckCategory(name string) (code int) {
	var cate Category
	db.Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCESS
}

// CreateCate 新增分类
func CreateCate(data *Category) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

// GetCateInfo 查询单个分类信息
func GetCateInfo(id int) (Category, int) {
    var cate Category
    if err := db.Where("id = ?", id).First(&cate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
            return cate, errmsg.ERROR_CATE_NOT_EXIST
        }
        return cate, errmsg.ERROR
    }
    return cate, errmsg.SUCCESS
}

// GetCate 查询分类列表
func GetCate(pageSize int, pageNum int) ([]Category, int64) {
    var cate []Category
    var total int64

    // 先计数
    if err := db.Model(&Category{}).Count(&total).Error; err != nil {
        return nil, 0
    }

    // 再分页查询
	if err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cate).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, 0
    }

    return cate, total
}

// EditCate 编辑分类信息
func EditCate(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err = db.Model(&cate).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteCate 删除分类
func DeleteCate(id int) int {
	var cate Category
	err = db.Where("id = ? ", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}