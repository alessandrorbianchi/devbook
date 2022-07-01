package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AppRabbitAddr      = ""
	InternalRun        = 0
	Queue              = ""
	StringConexaoBanco = ""
)

func Carregar() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	AppRabbitAddr = os.Getenv("APP_RABBIT_ADDR")
	InternalRun, err = strconv.Atoi(os.Getenv("INTERVAL_RUN"))
	if err != nil {
		InternalRun = 1
	}
	Queue = os.Getenv("QUEUE")

	StringConexaoBanco = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)
}
