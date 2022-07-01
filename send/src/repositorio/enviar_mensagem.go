package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"send/src/common"
	"send/src/config"
	"time"

	"github.com/streadway/amqp"
)

type Mensagem struct {
	ID             uint64 `json:"id,omitempty"`
	Mensagem       string `json:"mensagem,omitempty"`
	RemententeID   uint64 `json:"remetente_id,omitempty"`
	DestinatarioID uint64 `json:"destinatario_id,omitempty"`
}

type Send struct {
	channel *amqp.Channel
	db      *sql.DB
}

func NovoRepositorioRabbit(channel *amqp.Channel, db *sql.DB) *Send {
	return &Send{
		channel: channel,
		db:      db,
	}
}

func (s *Send) Enviar() error {
	fmt.Println("Executando enviar mensagem para o Rabbitmq")

	queue, err := s.channel.QueueDeclare(
		config.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		common.FailOnError(err, "Falha ao declarar uma fila")
	}

	var forever chan struct{}

	go func() {
		for {
			mensagens, err := s.BuscarMensagensNaoEnviadas()
			if err != nil {
				common.FailOnError(err, "Erro ao buscar mensagens não enviadas para a fila na base de dados")
			}

			for _, line := range mensagens {
				time.Sleep(time.Duration(config.InternalRun) * time.Second)

				payload, err := json.Marshal(line)
				if err != nil {
					common.FailOnError(err, "")
					continue
				}

				err = s.channel.Publish(
					"",         // exchange
					queue.Name, // routing key
					false,      // mandatory
					false,      // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        []byte(payload),
					})
				if err != nil {
					common.FailOnError(err, "Falha ao publicar uma mensagem")
				}

				if err := s.AtualizarDataEnviadoEm(line.ID); err != nil {
					common.FailOnError(err, "Erro ao atualizar campo EnviadoEm")
				}
			}
		}
	}()

	<-forever

	return nil
}

func (s *Send) BuscarMensagensNaoEnviadas() ([]Mensagem, error) {
	linhas, err := s.db.Query(
		`SELECT id, mensagem, remetente_id, destinatario_id
			FROM mensagens_enviadas 
			WHERE enviadoem is NULL`)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var mensagens []Mensagem

	for linhas.Next() {
		var mensagem Mensagem

		if err = linhas.Scan(
			&mensagem.ID,
			&mensagem.Mensagem,
			&mensagem.RemententeID,
			&mensagem.DestinatarioID,
		); err != nil {
			return nil, err
		}

		mensagens = append(mensagens, mensagem)
	}

	return mensagens, nil
}

func (s *Send) AtualizarDataEnviadoEm(ID uint64) error {
	statement, err := s.db.Prepare(
		"UPDATE mensagens_enviadas SET enviadoem = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(time.Now(), ID); err != nil {
		return err
	}

	return nil
}
