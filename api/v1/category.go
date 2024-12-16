package v1

import (
	"Go_blog/model"
	"Go_blog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 查询分类是否存在

// AddCategory 添加分类
func AddCategory(c *gin.Context) {
	var data model.Category
	c.ShouldBindJSON(&data)
	code:=model.CheckCategory(data.Name)
	if code==errmsg.SUCCESS {
		model.CreateCategory(&data)
	}

	c.JSON(
		http.StatusOK,gin.H{
			"status":code,
			"data":data,
			"message":errmsg.GetErrMsg(code),
		},
	)
}

// EditCategory 修改分类
func EditCategory(c *gin.Context) {
	var data model.Category
	id,_:=strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&data)
	code:=model.CheckCategory(data.Name)
	if code==errmsg.SUCCESS {
		model.EditCategory(id,&data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		c.Abort()
	}

	c.JSON(
		http.StatusOK,gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
		},
	)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id,_:=strconv.Atoi(c.Param("id"))
	code:=model.DeleteCategory(id)
	c.JSON(
		http.StatusOK,gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
		},
	)
}

// 查询分类信息
func GetCateInfo(c *gin.Context) {
	id,_:=strconv.Atoi(c.Param("id"))
	data,code:=model.GetCateInfo(id)
	c.JSON(
		http.StatusOK,gin.H{
			"status":code,
			"data":data,
			"message":errmsg.GetErrMsg(code),
		},
	)
}

// GetCategoryList 获取分类列表
func GetCategoryList(c *gin.Context) {
	pageSize,_:=strconv.Atoi(c.Query("pageSize"))
	pageNum,_:=strconv.Atoi(c.Query("pageNum"))

	switch{
	case pageSize>=100:
		pageSize=100
	case pageSize<=0:
		pageSize=10
	}
	if pageNum==0{
		pageNum=1
	}

	data,total:=model.GetCategory(pageSize,pageNum)
	code:=errmsg.SUCCESS
	c.JSON(
		http.StatusOK,gin.H{
			"status":code,
			"data":data,
			"total":total,
			"message":errmsg.GetErrMsg(code),
		},
	)
}
