package main

import (
	"send/src/banco"
	"send/src/common"
	"send/src/config"
	"send/src/repositorio"
	"send/src/server"
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

	enviarMensagem := repositorio.NovoRepositorioRabbit(
		channelRabbit,
		db,
	)

	if err := enviarMensagem.Enviar(); err != nil {
		common.FailOnError(err, "Erro ao enviar mensagem")
	}
}
