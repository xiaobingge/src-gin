package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"dbger/app/common"
	"dbger/app/models"
    "dbger/utils"
	"dbger/router"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var (
	configPath = pflag.StringP("config", "c", "", "dbger config file path.")
)
func init(){
	pflag.Parse()
	// load config include logger
	common.Init(*configPath)
	// Set gin mode.
	if "debug" == viper.GetString("runmode") {
		gin.SetMode(gin.DebugMode)
	} else if "test" == viper.GetString("runmode") {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
func main() {
	// init db
	models.DB.Init()
	defer models.DB.Close()
	// create the gin engine
	g := gin.New()
	// routes
	router.Load(g)
	_ = g.Run(viper.GetString("addr"))
	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up", err)
		}
		utils.Info("the router has been deployed successfully")
	}()
	// If open https, start listening https request
	if true == viper.GetBool("tls.https_open") {
		cert := viper.GetString("tls.cert")
		key := viper.GetString("tls.key")
		if cert != "" && key != "" {
			go func() {
				utils.Info("start to listening the incoming https requests", zap.String("port", viper.GetString("tls.addr")))
				utils.Info(http.ListenAndServeTLS("0.0.0.0:"+viper.GetString("tls.addr"), cert, key, g).Error())
			}()
		} else {
			utils.Errorf("cert and key can not be empty, failed to listen https port")
		}
	}
	utils.Info("start to listening the incoming http requests", zap.String("port", viper.GetString("addr")))
	utils.Info(http.ListenAndServe("0.0.0.0:"+viper.GetString("addr"), g).Error())

}

// pingServer pings the http server to make sure the service is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("ping_max_num"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("ping_url") + "/check/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		utils.Info("waiting for the router, retry in 1 second")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}