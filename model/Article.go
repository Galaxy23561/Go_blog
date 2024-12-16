package model

import (
	"Go_blog/utils/errmsg"
	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title        string `gorm:"type:varchar(100);not null" json:"title"`
	Cid          int    `gorm:"type:int;not null" json:"cid"`
	Desc         string `gorm:"type:varchar(200)" json:"desc"`
	Content      string `gorm:"type:longtext" json:"content"`
	Img          string `gorm:"type:varchar(100)" json:"img"`
	CommentCount int    `gorm:"type:int;not null;default:0" json:"comment_count"`
	ReadCount    int    `gorm:"type:int;not null;default:0" json:"read_count"`
}

// 新增文章
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询分类下的所有文章
func GetCateArt(id int,pageSize int,pageNum int) ([]Article,int,int64) {
	var cateArtList []Article
	var total int64

	err:=db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid=?",id).Find(&cateArtList).Error
	db.Model(&cateArtList).Where("cid=?",id).Count(&total)
	if err!=nil{
		return nil,errmsg.ERROR_CATE_NOT_EXIST,0
	}
	return cateArtList,errmsg.SUCCESS,total
}

// 查询单个文章
func GetArtInfo(id int) (Article,int) {
	var art Article
	err:=db.Where("id=?",id).Preload("Category").First(&art).Error
	db.Model(&art).Where("id=?",id).UpdateColumn("read_count",gorm.Expr("read_count + ?",1))
	if err !=nil{
		return art,errmsg.ERROR_ART_NOT_EXIST
	}
	return art,errmsg.SUCCESS
}

// 查询文章列表
func GetArt(pageSize int,pageNum int) ([]Article,int,int64) {
	var artList []Article
	var total int64
	err:=db.Select("article.id,title,img,created_at,updated_at,`desc`,comment_count,read_count,category.name").Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_AT DESC").Joins("Category").Find(&artList).Error
	db.Model(&artList).Count(&total)
	if err!=nil{
		return nil,errmsg.ERROR,0
	}
	return artList,errmsg.SUCCESS,total
}

// 编辑文章
func EditArticle(id int,data *Article) int {
	var art Article
	var maps=make(map[string]interface{})
	maps["title"]=data.Title
	maps["cid"]=data.Cid
	maps["desc"]=data.Desc
	maps["content"]=data.Content
	maps["img"]=data.Img

	err:=db.Model(&art).Where("id=?",id).Updates(&maps).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 搜索文章标题
func SearchArt(title string,pageSize int,pageNum int)([]Article,int,int64) {
	var artList []Article
	var total int64
	err:=db.Select("article.id,title,img,created_at,updated_at,`desc`,comment_count,read_count,category.name").Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_AT DESC").Joins("Category").Where("title LIKE ?","%"+title+"%").Find(&artList).Error
	db.Model(&artList).Where("title LIKE ?","%"+title+"%").Count(&total)
	if err!=nil{
		return nil,errmsg.ERROR,0
	}
	return artList,errmsg.SUCCESS,total
}

// 删除文章
func DeleteArticle(id int) int {
	var art Article
	err:=db.Where("id=?",id).Delete(&art).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}