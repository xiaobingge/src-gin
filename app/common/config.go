package common

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xiaobingge/dbger/app/utils"
	"go.uber.org/zap"
	"strings"
)

// Config config struct
type Config struct {
	Name string
}

// Init sets all configs using config file setting.
func Init(cfg string) {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		panic(fmt.Errorf("fatal error init configuration: %s", err))
	}

	// 初始化日志包
	c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	utils.Info("configuration load succeeded", zap.String("config file", viper.ConfigFileUsed()))
}

// initConfig int config
func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，解析指定配置文件
	} else {
		viper.AddConfigPath("./configs") // 如果没有指定配置文件，则解析默认配置文件
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("DBGER")  //读取环境变量的前缀为DBGER

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

// initLog init log config 初始化日志配置
func (c *Config) initLog() {
	utils.InitLogger(viper.GetString("runmode"))
	utils.Info("zap logger construction succeeded")
}

// watchConfig watch config if has been changed 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		utils.Infof("config file changed: %s", e.Name)
	})
}