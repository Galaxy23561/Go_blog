package main

import (
	"Go_blog/model"
	"Go_blog/routes"
)

func main(){
	model.InitDb()
	routes.InitRouter()
}