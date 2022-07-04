package modelos

import "time"

type Mensagem struct {
	ID               uint64    `json:"id"`
	Mensagem         string    `json:"mensagem"`
	RemetenteID      uint64    `json:"remetente_id"`
	RemetenteNick    string    `json:"remetente_nick"`
	DestinatarioID   uint64    `json:"destinatario_id"`
	DestinatarioNick string    `json:"destinatario_nick"`
	CriadoEm         time.Time `json:"criadoem"`
	EnviadoEm        string    `json:"enviadoem"`
	RecebidoEm       string    `json:"recebidoem"`
}
