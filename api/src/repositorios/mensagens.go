package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type Mensagens struct {
	db *sql.DB
}

func NovoRepositorioDeMensagens(db *sql.DB) *Mensagens {
	return &Mensagens{db}
}

func (m Mensagens) Enviar(mensagem modelos.MensagemEnviada) (uint64, error) {
	statement, err := m.db.Prepare("INSERT INTO mensagens_enviadas (mensagem, remetente_id, destinatario_id, codigo_seguranca) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(mensagem.Mensagem, mensagem.RemetenteID, mensagem.DestinatarioID, mensagem.CodigoSeguranca)
	if err != nil {
		return 0, err
	}

	ultimoIDInserido, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoIDInserido), nil
}

func (m Mensagens) BuscarAgrupado(ID uint64) ([]modelos.MensagemEnviada, error) {
	linhas, err := m.db.Query(
		`SELECT 
			me.id, 
			me.mensagem, 
			me.remetente_id, 
			ur.nick, 
			me.destinatario_id, 
			ud.nick, 
			me.codigo_seguranca, 
			me.criadoem, 
			COALESCE(me.enviadoem, 'NULL'), 
			COALESCE(mr.recebidoem, 'NULL')
		FROM mensagens_enviadas me  
		INNER JOIN usuarios ur ON ur.id = me.remetente_id 
		INNER JOIN usuarios ud ON ud.id = me.destinatario_id 
		LEFT JOIN mensagens_recebidas mr ON me.id = mr.mensagem_enviada_id 
		WHERE me.id = (
			SELECT 
				MAX(me2.id) 
			FROM mensagens_enviadas me2
			WHERE 
				me2.destinatario_id = ? 
				OR me2.remetente_id  = ? 
			)`,
		ID, ID,
	)

	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var mensagens []modelos.MensagemEnviada

	for linhas.Next() {
		var mensagem modelos.MensagemEnviada

		if err = linhas.Scan(
			&mensagem.ID,
			&mensagem.Mensagem,
			&mensagem.RemetenteID,
			&mensagem.RemetenteNick,
			&mensagem.DestinatarioID,
			&mensagem.DestinatarioNick,
			&mensagem.CodigoSeguranca,
			&mensagem.CriadoEm,
			&mensagem.EnviadoEm,
			&mensagem.RecebidoEm,
		); err != nil {
			return nil, err
		}

		mensagens = append(mensagens, mensagem)
	}

	return mensagens, nil
}

func (m Mensagens) BuscarPorUsuario(remetenteID, destinatarioID uint64) ([]modelos.MensagemEnviada, error) {
	linha, err := m.db.Query(
		`SELECT me.id, me.mensagem, me.remetente_id, ur.nick, me.destinatario_id, ud.nick, me.criadoem, COALESCE(me.enviadoem, 'NULL'), COALESCE(mr.recebidoem, 'NULL')
			FROM mensagens_enviadas me  
			INNER JOIN usuarios ur ON ur.id = me.remetente_id 
			INNER JOIN usuarios ud ON ud.id = me.destinatario_id 
			LEFT JOIN mensagens_recebidas mr ON me.id = mr.mensagem_enviada_id 
		WHERE me.destinatario_id = ? and me.remetente_id  = ?
		OR me.destinatario_id = ? and me.remetente_id  = ?
		AND mr.recebidoem IS NOT NULL
		ORDER BY me.id DESC 
		`,
		destinatarioID, remetenteID, remetenteID, destinatarioID,
	)
	if err != nil {
		return nil, err
	}
	defer linha.Close()

	var mensagens []modelos.MensagemEnviada

	for linha.Next() {
		var mensagem modelos.MensagemEnviada

		if err = linha.Scan(
			&mensagem.ID,
			&mensagem.Mensagem,
			&mensagem.RemetenteID,
			&mensagem.RemetenteNick,
			&mensagem.DestinatarioID,
			&mensagem.DestinatarioNick,
			&mensagem.CriadoEm,
			&mensagem.EnviadoEm,
			&mensagem.RecebidoEm,
		); err != nil {
			return nil, err
		}

		mensagens = append(mensagens, mensagem)
	}

	return mensagens, nil
}
