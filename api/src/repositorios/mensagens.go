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
	statement, err := m.db.Prepare("INSERT INTO mensagens_enviadas (mensagem, remetente_id, destinatario_id) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(mensagem.Mensagem, mensagem.RemententeID, mensagem.DestinatarioID)
	if err != nil {
		return 0, err
	}

	ultimoIDInserido, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoIDInserido), nil
}

func (m Mensagens) Buscar(ID uint64) ([]modelos.MensagemEnviada, error) {
	linhas, err := m.db.Query(
		`SELECT me.id, me.mensagem, me.remetente_id, me.destinatario_id , me.criadoem, COALESCE(me.enviadoem, 'NULL')
		 	FROM mensagens_enviadas me
		 	LEFT JOIN mensagens_recebidas mr ON me.id = mr.mensagem_enviada_id
		 WHERE me.remetente_id = ? or me.destinatario_id = ?
		order by 1 desc`,
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
			&mensagem.RemententeID,
			&mensagem.DestinatarioID,
			&mensagem.CriadoEm,
			&mensagem.EnviadoEm,
		); err != nil {
			return nil, err
		}

		mensagens = append(mensagens, mensagem)
	}

	return mensagens, nil
}

func (m Mensagens) BuscarPorUsuario(usuarioID uint64) ([]modelos.MensagemEnviada, error) {
	linha, err := m.db.Query(
		`SELECT me.id, me.mensagem, me.remetente_id, me.destinatario_id , me.criadoem, COALESCE(me.enviadoem, 'null')
			FROM mensagens_enviadas me
			INNER JOIN usuarios u ON u.id = me.remetente_id
			WHERE me.remetente_id = ?`,
		usuarioID,
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
			&mensagem.RemententeID,
			&mensagem.DestinatarioID,
			&mensagem.CriadoEm,
			&mensagem.EnviadoEm,
		); err != nil {
			return nil, err
		}

		mensagens = append(mensagens, mensagem)
	}

	return mensagens, nil
}

// func (m Mensagens) BuscarPorID(mensagemID uint64) (modelos.MensagemEnviada, error) {
// 	linha, err := m.db.Query(
// 		`SELECT me.id, me.mensagem, me.remetente_id, me.destinatario_id , me.criadoem, COALESCE(me.enviadoem, 'null')
// 		 	FROM mensagens_enviadas me
// 		 	LEFT JOIN usuarios u ON u.id = me.remetente_id
// 		 WHERE me.remetente_id = ?`,
// 		mensagemID,
// 	)
// 	if err != nil {
// 		return modelos.MensagemEnviada{}, err
// 	}
// 	defer linha.Close()

// 	var mensagem modelos.MensagemEnviada

// 	if linha.Next() {
// 		if err = linha.Scan(
// 			&mensagem.ID,
// 			&mensagem.Mensagem,
// 			&mensagem.RemententeID,
// 			&mensagem.DestinatarioID,
// 			&mensagem.CriadoEm,
// 			&mensagem.EnviadoEm,
// 		); err != nil {
// 			return modelos.MensagemEnviada{}, err
// 		}
// 	}

// 	return mensagem, nil
// }

// func (p Publicacoes) Deletar(publicacaoID uint64) error {
// 	statement, err := p.db.Prepare("DELETE FROM publicacoes WHERE id = ?")

// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	if _, err = statement.Exec(publicacaoID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (p Publicacoes) Descurtir(publicacaoID uint64) error {
// 	statement, err := p.db.Prepare(
// 		`UPDATE publicacoes SET curtidas =
// 			CASE
// 				WHEN curtidas > 0 THEN curtidas - 1
// 				ELSE 0
// 			END
// 			WHERE id = ?`,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	if _, err = statement.Exec(publicacaoID); err != nil {
// 		return err
// 	}

// 	return nil
// }
