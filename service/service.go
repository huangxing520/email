package service

import (
	constant2 "emailTest/constant"
	db "emailTest/db/database"
	"emailTest/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/jordan-wright/email"
	"github.com/robfig/cron/v3"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

var (
	ListenerListInstance *ListenerList
	once                 sync.Once
)

type ListenerList struct {
	list map[string]Listener
	c    *cron.Cron
	ListenerHandle
}
type Listener struct {
	birthdayDate string
	currentCron  cron.EntryID
	preCron      cron.EntryID
}
type ListenerHandle interface {
	Add(name string, birthday string) error
	Delete(name string)
}

func InitService() {
	listeners := GetInstance()
	listeners.c.Start()
	Server()
}
func (t *ListenerList) Add(name string, birthdayDate string) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error(err)
		}
	}()
	currentCron, err := t.c.AddFunc(BirthdayToCorn(birthdayDate), func() {
		SendEmail(name)
	})
	if err != nil {
		fmt.Println(err)
	}
	preCron, err := t.c.AddFunc(PreBirthdayToCorn(birthdayDate), func() {
		fmt.Println("开始发送")
		preSendEmail(name, birthdayDate)
	})
	log.Logger.Info(name, "请求创建订阅，日期是", birthdayDate)
	t.list[name] = Listener{
		birthdayDate: birthdayDate,
		currentCron:  currentCron,
		preCron:      preCron,
	}
	db.SqliteDb.AddPerson(name, birthdayDate)
	db.AddJson(name, birthdayDate)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}
func Server() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	emailGroup := r.Group("/email")
	{
		emailGroup.GET("/add", func(c *gin.Context) {
			listeners := GetInstance()
			listeners.Add(c.Query("name"), c.Query("birthday"))
			c.String(http.StatusOK, "add_ok")
		})
		emailGroup.GET("/delete", func(c *gin.Context) {
			listeners := GetInstance()
			listeners.Delete(c.Query("name"))
			c.String(http.StatusOK, "delete_ok")
		})
		emailGroup.GET("/list", func(c *gin.Context) {
			c.String(http.StatusOK, db.List())
		})
	}

	err := r.Run(":8011")
	if err != nil {
		log.Logger.Error(err)
		return
	}
}
func (t *ListenerList) Delete(name string) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error(err)
		}
	}()
	t.c.Remove(t.list[name].currentCron)
	t.c.Remove(t.list[name].preCron)
	delete(t.list, name)
	db.SqliteDb.Delete(name)
	db.DeleteJson(name)
	log.Logger.Info(name, "请求取消订阅")

}

func SendEmail(name string) {
	from := constant2.FromEmail
	To := constant2.ToEmail
	passwd := constant2.EmailPasswd
	e := &email.Email{
		To:      []string{To},
		From:    from,
		Subject: "【提醒】【生日】" + name,
		Text:    []byte(name + "今天生日！"),
		HTML:    []byte(""),
		Headers: textproto.MIMEHeader{},
	}
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", from, passwd, "smtp.qq.com"))
	if err != nil {
		log.Logger.Error(err)
		return
	}
	log.Logger.Info("发送给", name, "邮件成功，时间是：", time.Now())
}
func preSendEmail(name string, date string) {
	from := constant2.FromEmail
	To := constant2.ToEmail
	passwd := constant2.EmailPasswd
	e := &email.Email{
		To:      []string{To},
		From:    from,
		Subject: "【预提醒】【生日】" + name,
		Text:    []byte(name + "生日即将到来，" + "时间: " + date),
		HTML:    []byte(""),
		Headers: textproto.MIMEHeader{},
	}
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", from, passwd, "smtp.qq.com"))
	if err != nil {
		log.Logger.Error(err)
		return
	}
	log.Logger.Info("预发送给", name, "邮件成功，发送时间是：", time.Now(), "生日日期是：", date)
}

// BirthdayToCorn birthday:5 20/**
func BirthdayToCorn(birthday string) string {
	splits := strings.Split(birthday, "_")
	return "0 0 " + splits[1] + " " + splits[0] + " *"
}
func PreBirthdayToCorn(birthday string) string {
	birth := strings.Replace(birthday, "_", "-", -1)
	preDate := carbon.Parse(fmt.Sprintf("%v-%v", carbon.Now().Year(), birth)).AddDays(-5).ToDateTimeString()
	fmt.Println(carbon.Parse(fmt.Sprintf("%v-%v", carbon.Now().Year(), birth)).AddDays(-5).ToDateTimeString())
	month := strings.Split(strings.Split(preDate, " ")[0], "-")[1]
	day := strings.Split(strings.Split(preDate, " ")[0], "-")[2]
	if month == "02" && day == "29" {
		day = "28"
	}
	return "0 0 " + day + " " + month + " *"

}

func GetInstance() *ListenerList {
	once.Do(
		func() {
			ListenerListInstance = &ListenerList{list: make(map[string]Listener, 10),
				c: cron.New()}
		})

	return ListenerListInstance
}
