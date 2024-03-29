package config

import (
	"flag"
	"log"

	"github.com/notnull-co/cfg"
)

var (
	Cfg Config
)

type Config struct {
	Token struct {
		Key string `cfg:"key"`
	} `cfg:"token"`
	Server struct {
		Port string `cfg:"port"`
	} `cfg:"server"`
	DB struct {
		ConnectionString string `cfg:"connection_string"`
	} `cfg:"db"`
	SQS struct {
		PaymentPendingQueue   string `cfg:"payment_pending_queue"`
		PaymentPayedQueue     string `cfg:"payment_payed_queue"`
		PaymentCancelledQueue string `cfg:"payment_cancelled_queue"`
		Region                string `cfg:"region"`
	} `cfg:"sqs"`
}

func ParseFromFlags() {
	var configDir string

	flag.StringVar(&configDir, "config-dir", "../../internal/config/", "Configuration file directory")
	flag.Parse()

	parse(configDir)
}

func parse(dirs ...string) {
	if err := cfg.Load(&Cfg,
		cfg.Dirs(dirs...),
		cfg.UseEnv("app"),
	); err != nil {
		log.Panic(err)
	}
}
func Get() Config {
	return Cfg
}
