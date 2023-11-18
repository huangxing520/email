package db

import (
	co "emailTest/constant"
	log "emailTest/logger"
	"encoding/json"
	"os"
)

type JsonUser struct {
	Name     string `json:"Name"`
	Birthday string `json:"Birthday"`
}

var userList []JsonUser

func InitJson(path string) []JsonUser {
	jsonfile, err := os.Open(path)
	if err != nil {
		log.Logger.Error("打开json文件失败")
		return nil
	}
	defer jsonfile.Close()

	decoder := json.NewDecoder(jsonfile) // 创建 json 解码器
	err = decoder.Decode(&userList)
	return userList
}
func EncodeJson(path string) error {
	jsonfile, err := os.Create(path)
	if err != nil {
		log.Logger.Error("打开json文件失败")
		return err
	}
	defer jsonfile.Close()
	err = json.NewEncoder(jsonfile).Encode(&userList)
	if err != nil {
		return err
	}
	return nil
}
func AddJson(name string, birthday string) {
	var index = len(userList)
	for i, v := range userList {
		if v.Name == name {
			index = i
		}
	}
	if index != len(userList) {
		return
	}
	userList = append(userList, JsonUser{
		Name:     name,
		Birthday: birthday,
	})
	err := EncodeJson(co.JsonPath)
	if err != nil {
		return
	}
}
func DeleteJson(name string) {
	var index = len(userList)
	for i, v := range userList {
		if v.Name == name {
			index = i
		}
	}
	if index == len(userList) {
		return
	}
	userList = append(userList[:index], userList[index+1:]...)
	err := EncodeJson(co.JsonPath)
	if err != nil {
		return
	}
}
func List() string {
	jsonData, _ := json.Marshal(userList)
	return string(jsonData)
}
