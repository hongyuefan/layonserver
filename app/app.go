package app

import (
	"fmt"
	"layonserver/api"
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
	}
	if len(g_ConfigData.MongoConn) > 0 {
		mongodb.DitalMongo(g_ConfigData.MongoConn)
	}

	app.handlers = api.NewHandlers()

	router := gin.Default()

	router.Handle("GET", "health", func(c *gin.Context) { c.String(200, "ok") })

	v1 := router.Group("/v1") /*.Use(app.handlers.UserAuthrize)*/
	{
		v1.POST("/post", app.handlers.HandlerPost)
		v1.GET("/get", app.handlers.HandlerGet)
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
