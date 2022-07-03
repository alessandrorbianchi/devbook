package modelos

import (
	"errors"
	"strings"
)

type MensagemEnviada struct {
	ID               uint64 `json:"id,omitempty"`
	Mensagem         string `json:"mensagem,omitempty"`
	RemententeID     uint64 `json:"rementente_id,omitempty"`
	RemententeNick   string `json:"remetente_nick,omitempty"`
	DestinatarioID   uint64 `json:"destinatario_id,omitempty"`
	DestinatarioNick string `json:"destinatario_nick,omitempty"`
	CodigoSeguranca  uint64 `json:"codigo_seguranca,omitempty"`
	CriadoEm         string `json:"criadoem,omitempty"`
	EnviadoEm        string `json:"enviadoem,omitempty"`
	RecebidoEm       string `json:"recebido_em,omitempty"`
}

func (m *MensagemEnviada) Preparar() error {
	if err := m.validar(); err != nil {
		return err
	}

	if err := m.formatar(); err != nil {
		return err
	}

	return nil
}

func (m *MensagemEnviada) validar() error {
	if m.Mensagem == "" {
		return errors.New("a mensagem é obrigatória e não pode estar em branco")
	}

	return nil
}

func (m *MensagemEnviada) formatar() error {
	m.Mensagem = strings.TrimSpace(m.Mensagem)

	return nil
}
