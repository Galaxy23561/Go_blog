package model

import (
	"Go_blog/utils/errmsg"
	"errors"
	"gorm.io/gorm"
	"log"
)

type Article struct {
	gorm.Model
	Title        string   `gorm:"type:varchar(100);not null" json:"title"`
	Cid          uint     `json:"cid"` // 确保类型匹配 Category 的 ID 类型
	Category     Category `gorm:"foreignKey:Cid;references:ID"`
	Desc         string   `gorm:"type:varchar(200)" json:"desc"`
	Content      string   `gorm:"type:longtext" json:"content"`
	Img          string   `gorm:"type:varchar(100)" json:"img"`
	CommentCount int      `gorm:"type:int;not null;default:0" json:"comment_count"`
	ReadCount    int      `gorm:"type:int;not null;default:0" json:"read_count"`
	Info         string   `gorm:"type:text;not null" json:"info"` // 添加 info 字段
}

// CreateArt 新增文章
/*func CreateArt(data *Article) int {
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}*/
// CreateArticle 创建文章的方法
func (a *Article) CreateArticle() (int, string) {
	db := GetDB()

	// 验证标题是否为空
	if a.Title == "" {
		return errmsg.ERROR_TITLE_NOT_EXIST, "标题不能为空"
	}

	// 检查分类是否存在
	var category Category
	if err := db.Where("id = ?", a.Cid).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR_CATE_NOT_EXIST, "分类不存在"
		}
		log.Printf("数据库错误: %v", err)
		return errmsg.ERROR, "数据库错误"
	}

	// 使用事务来保证数据一致性
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建文章
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		log.Printf("创建文章失败: %v", err)
		return errmsg.ERROR, "创建文章失败"
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("事务提交失败: %v", err)
		return errmsg.ERROR, "事务提交失败"
	}

	log.Printf("文章创建成功: %d", a.ID)
	return errmsg.SUCCESS, "文章创建成功"
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64

	// 先计数再分页查询
	if err := db.Model(&Article{}).Where("cid = ?", id).Count(&total).Error; err != nil {
		return nil, errmsg.ERROR, 0
	}

	if err := db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"cid = ?", id).Find(&cateArtList).Error; err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}

	return cateArtList, errmsg.SUCCESS, total
}

// GetArtInfo 查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id = ?", id).Preload("Category").First(&art).Error; err != nil {
		tx.Rollback()
		return art, errmsg.ERROR_ART_NOT_EXIST
	}

	if err := tx.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return art, errmsg.ERROR
	}

	tx.Commit()
	return art, errmsg.SUCCESS
}

// GetArt 查询文章列表
func GetArt(pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var total int64

	// 计数
	if err := db.Model(&Article{}).Count(&total).Error; err != nil {
		log.Printf("Failed to count articles: %v", err)
		return nil, errmsg.ERROR, 0
	}

	// 分页查询
	query := db.Select("article.id, article.title, article.img, article.created_at, article.updated_at, article.desc, article.comment_count, article.read_count, category.name as category_name").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Order("article.created_at DESC").
		Joins("LEFT JOIN category ON article.cid = category.id")

	if err := query.Find(&articleList).Error; err != nil {
		log.Printf("Failed to fetch articles: %v", err)
		return nil, errmsg.ERROR, 0
	}

	return articleList, errmsg.SUCCESS, total
}

// SearchArticle 搜索文章标题
func SearchArticle(title string, pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var total int64

	// 计数
	if err := db.Model(&Article{}).Where("title LIKE ?", "%"+title+"%").Count(&total).Error; err != nil {
		return nil, errmsg.ERROR, 0
	}

	// 分页查询
	if err := db.Select("id, title, img, created_at, updated_at, `desc`, comment_count, read_count, Category.name").
		Where("title LIKE ?", "%"+title+"%").
		Order("created_at DESC").
		Joins("LEFT JOIN categories ON articles.cid = categories.id").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&articleList).Error; err != nil {
		return nil, errmsg.ERROR, 0
	}

	return articleList, errmsg.SUCCESS, total
}

// EditArt 编辑文章
func EditArt(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err = db.Model(&art).Where("id = ? ", id).Updates(&maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteArt 删除文章
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
