package main

import (
	"emailTest/constant"
	db "emailTest/db/database"
	"emailTest/service"
)

func main() {
	db.InitDb("./email.db")
	userList := db.InitJson(constant.JsonPath)
	listeners := service.GetInstance()
	for _, v := range userList {
		listeners.Add(v.Name, v.Birthday)
	}
	service.InitService()

}
