package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"receive/src/common"
	"receive/src/config"
	"time"

	"github.com/streadway/amqp"
)

type Receive struct {
	channel *amqp.Channel
	db      *sql.DB
}

type Mensagem struct {
	ID             uint64 `json:"id,omitempty"`
	Mensagem       string `json:"mensagem,omitempty"`
	RemententeID   uint64 `json:"remetente_id,omitempty"`
	DestinatarioID uint64 `json:"destinatario_id,omitempty"`
}

func NovoRepositorioRabbit(channel *amqp.Channel, db *sql.DB) *Receive {
	return &Receive{
		channel: channel,
		db:      db,
	}
}

func (r *Receive) Receber() error {
	fmt.Println("Executando receber mensagem do Rabbitmq")

	queue, err := r.channel.QueueDeclare(
		config.Queue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		common.FailOnError(err, "Falha ao declarar uma fila")
	}

	msgs, err := r.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		common.FailOnError(err, "Falha ao registrar um consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var mensagem Mensagem
			err = json.Unmarshal(d.Body, &mensagem)

			if err := r.Inserir(mensagem); err != nil {
				log.Printf("Erro ao inserir dados na tabela: %s - mensagem %v", err, mensagem)
				continue
			}
			time.Sleep(time.Duration(config.InternalRun) * time.Second)
		}
	}()

	log.Printf(" [*] Aguardando mensagens. Para sair pressione CTRL+C")
	<-forever

	return nil
}

func (r *Receive) Inserir(mensagem Mensagem) error {
	statement, err := r.db.Prepare(
		`INSERT INTO mensagens_recebidas 
			(mensagem_enviada_id, mensagem, remetente_id, destinatario_id) 
		VALUES(?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(mensagem.ID, mensagem.Mensagem, mensagem.RemententeID, mensagem.DestinatarioID)
	if err != nil {
		return err
	}

	return nil
}
