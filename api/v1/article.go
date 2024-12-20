package v1

import (
	"github.com/gin-gonic/gin"
	"Go_blog/model"
	"Go_blog/utils/errmsg"
	"net/http"
	"strconv"
)

// AddArticle 添加文章
/*func AddArticle(c *gin.Context) {
	var data model.Article
	_ = c.ShouldBindJSON(&data)

	code := model.CreateArt(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}*/
// AddArticle 添加文章的API处理函数
func AddArticle(c *gin.Context) {
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR_INVALID_PARAM,
			"message": errmsg.GetErrMsg(errmsg.ERROR_INVALID_PARAM),
		})
		return
	}

	// 调用模型中的CreateArticle方法
	status, message := data.CreateArticle()
	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"data":    data,
		"message": message,
	})

	// 如果创建失败，返回错误状态码
	if status != errmsg.SUCCESS {
		c.Status(http.StatusInternalServerError)
	}
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	id, _ := strconv.Atoi(c.Param("id"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, code, total := model.GetCateArt(id, pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArtInfo 查询单个文章信息
func GetArtInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetArtInfo(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArt 查询文章列表
func GetArt(c *gin.Context) {
	pageSizeStr := c.Query("pagesize")
	pageNumStr := c.Query("pagenum")
	title := c.Query("title")

	// 默认分页参数
	pageSize := 10
	pageNum := 1

	// 解析分页参数
	if pageSizeInt, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeInt > 0 {
		pageSize = pageSizeInt
		if pageSize > 100 {
			pageSize = 100
		}
	}

	if pageNumInt, err := strconv.Atoi(pageNumStr); err == nil && pageNumInt > 0 {
		pageNum = pageNumInt
	}

	// 查询文章列表或搜索文章
	var data []model.Article
	var code int
	var total int64
	var message string

	if len(title) == 0 {
		data, code, total = model.GetArt(pageSize, pageNum)
	} else {
		data, code, total = model.SearchArticle(title, pageSize, pageNum)
	}

	message = errmsg.GetErrMsg(code)

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": message,
	})

	// 如果状态码不是成功，则设置相应的 HTTP 状态码
	if code != errmsg.SUCCESS {
		c.Status(http.StatusInternalServerError)
	}
}

// EditArt 编辑文章
func EditArt(c *gin.Context) {
	var data model.Article
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.EditArt(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteArt 删除文章
func DeleteArt(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteArt(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}