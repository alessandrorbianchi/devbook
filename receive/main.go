package main

import (
	"receive/src/banco"
	"receive/src/common"
	"receive/src/config"
	"receive/src/repositorio"
	"receive/src/server"
)

func main() {
	config.Carregar()

	rabbit, err := server.ConectarNoRabbit(config.AppRabbitAddr)
	if err != nil {
		common.FailOnError(err, "Falha ao conectar no RabbitMQ")
	}
	defer rabbit.Close()

	channelRabbit, err := rabbit.Channel()
	if err != nil {
		common.FailOnError(err, "Falha ao abrir um canal")
	}
	defer channelRabbit.Close()

	db, err := banco.Conectar()
	if err != nil {
		common.FailOnError(err, "Erro ao conectar com o banco de dados")
		return
	}
	defer db.Close()

	receberMensagem := repositorio.NovoRepositorioRabbit(
		channelRabbit,
		db,
	)

	if err := receberMensagem.Receber(); err != nil {
		common.FailOnError(err, "Erro ao receber mensagem")
	}
}
