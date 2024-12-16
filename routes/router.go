package routes

import (
	"github.com/gin-gonic/gin"
	"Go_blog/utils"
	v1 "Go_blog/api/v1"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.Default()

	routerV1 := r.Group("api/v1")
	{
		/*
		router.GET("hello",func (c *gin.Context){
			c.JSON(200,gin.H{
				"msg":"ok",
			})
		})
		*/
		// 用户模块的路由接口
		routerV1.POST("user/add", v1.AddUser)
		routerV1.GET("users", v1.GetUserList)
		routerV1.PUT("user/:id", v1.EditUser)
		routerV1.DELETE("user/:id", v1.DeleteUser)


		// 分类模块的路由接口
		routerV1.POST("category/add",v1.AddCategory)
		routerV1.GET("categories",v1.GetCategoryList)
		routerV1.PUT("category/:id",v1.EditCategory)
		routerV1.DELETE("category/:id",v1.DeleteCategory)
		// 文章模块的路由接口
		routerV1.POST("article/add",v1.AddArticle)
		routerV1.GET("articles",v1.GetArticleList)
		routerV1.GET("article/info/:id",v1.GetSingleArticle)
		routerV1.GET("category/articles/:id",v1.GetCateArticle)
		routerV1.PUT("article/:id",v1.EditArticle)
		routerV1.DELETE("article/:id",v1.DeleteArticle)
		// 上传文件
	}
	r.Run(utils.HttpPort)
	//r.Run("0.0.0.0:3000")
}