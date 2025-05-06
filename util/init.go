package util

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	logrus_logstash "github.com/sima-land/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logFile = "/tmp/payso-external-api.log"
var logLevel = "DEBUG"

func Init() {
	log.Infof("==-- Init Env --==")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	log.Infof("ENV_CONFIG => %s", viper.GetString("ENV_CONFIG"))
	log.Infof("SECRET => %s", viper.GetString("SECRET_NAME"))
	//Get data from AWS Secret Manager
	if strings.Compare("LOCAL", viper.GetString("ENV_CONFIG")) != 0 {
		log.Info("==-- GET ENV From Server --==")
		envval, err := InitEnv(viper.GetString("SECRET_NAME"), viper.GetString("AWS_REGION"))

		if err != nil {
			panic(err)
		}

		for key, element := range envval {
			//fmt.Println("==>", key, "-", element)
			viper.Set(key, element)
		}
	} else {
		log.Info("==-- GET ENV From LOCAL --==")
	}

	localIP, err := LocalIP()

	if err != nil {
		viper.Set("MACHINE", "127.0.0.1")
	}

	fmt.Println("local ip :", localIP)
	viper.Set("MACHINE", localIP)

	logLevel = viper.GetString("LOG_LEVEL")
	println("log:" + logLevel)

	file, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", logFile, "%Y-%m-%d"),
		rotatelogs.WithLinkName(logFile+".link"),
		rotatelogs.WithMaxAge(time.Hour*24*10),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	mw := io.MultiWriter(os.Stdout, file)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        true,
		NoColors:        false,
		FieldsOrder:     []string{"component", "function"},
	})

	log.SetOutput(mw)

	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

	//Log Stash
	hook, err := logrus_logstash.NewHook("tcp", viper.GetString("LOGSTASH"), "simple-noti")
	if err != nil {
		log.Error(err)
	} else {
		log.Info("==-- Add Log Stash --==")
		log.AddHook(hook)
	}

}
