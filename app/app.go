package app

import (
	"fmt"
	"layonserver/api"
	"layonserver/models"
	"layonserver/util/config"
	"layonserver/util/log"
	"layonserver/util/mongodb"
	"layonserver/util/qmsql"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

const MasterName = "layonserver"

type ConfigData struct {
	Port      string
	LogDir    string
	IsHttps   bool
	Cert      string
	Key       string
	SqlConn   string
	MongoConn string
}

type App struct {
	handlers *api.Handlers
}

var g_ConfigData *ConfigData

func OnInitFlag(c *config.Config) (err error) {

	g_ConfigData = new(ConfigData)
	g_ConfigData.Port = c.GetString("port")
	g_ConfigData.LogDir = c.GetString("logdir")
	g_ConfigData.IsHttps = c.GetBool("https")
	g_ConfigData.Cert = c.GetString("https_cert")
	g_ConfigData.Key = c.GetString("https_key")
	g_ConfigData.SqlConn = c.GetString("sql_url")
	g_ConfigData.MongoConn = c.GetString("mongo_url")

	if "" == g_ConfigData.Port || "" == g_ConfigData.LogDir {
		return fmt.Errorf("config not right")
	}
	return

}

func (app *App) OnStart(c *config.Config) error {

	if err := OnInitFlag(c); err != nil {
		return err
	}

	if _, err := log.NewLog(g_ConfigData.LogDir, MasterName, 0); err != nil {
		return err
	}

	if len(g_ConfigData.SqlConn) > 0 {
		if err := qmsql.InitMysql(g_ConfigData.SqlConn); err != nil {
			panic(err)
		}
		qmsql.DEFAULTDB.AutoMigrate(models.Users{}, models.Devices{})
	}
	if len(g_ConfigData.MongoConn) > 0 {
		mongodb.DitalMongo(g_ConfigData.MongoConn)
	}

	app.handlers = api.NewHandlers()

	router := gin.Default()

	router.Handle("GET", "health", func(c *gin.Context) { c.String(200, "ok") })

	lg := router.Group("/lg").Use(app.handlers.UserAuthrize)
	{
		lg.POST("/device/del", app.handlers.HandlerDelDevice)
		lg.POST("/device/list", app.handlers.HandlerGetDevices)
		lg.POST("/device/add", app.handlers.HandlerAddDevice)
		lg.POST("/users/list", app.handlers.HandlerGetUsers)
		lg.POST("/users/info", app.handlers.HandlerGetUser)
		lg.POST("/users/listbyfather", app.handlers.HandlerGetUsersByFather)
	}

	rg := router.Group("/rg")
	{
		rg.GET("/verifycode", app.handlers.HandlerGetVerifyCode)
		rg.POST("/login", app.handlers.HandlerLogin)
		rg.POST("/regist", app.handlers.HandlerRegist)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	if g_ConfigData.IsHttps {
		http.ListenAndServeTLS(":"+g_ConfigData.Port, g_ConfigData.Cert, g_ConfigData.Key, router)
	} else {
		http.ListenAndServe(":"+g_ConfigData.Port, router)
	}
	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	fmt.Println("server shutdown")
}
