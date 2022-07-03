package modelos

import "time"

type Mensagem struct {
	ID               uint64    `json:"id"`
	Mensagem         string    `json:"mensagem"`
	RemententeID     uint64    `json:"rementente_id"`
	RemententeNick   string    `json:"remetente_nick"`
	DestinatarioID   uint64    `json:"destinatario_id"`
	DestinatarioNick string    `json:"destinatario_nick"`
	CriadoEm         time.Time `json:"criadoem"`
	EnviadoEm        string    `json:"enviadoem,omitempty"`
	RecebidoEm       string    `json:"recebido_em,omitempty"`
}
