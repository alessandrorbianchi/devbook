package modelos

import (
	"errors"
	"strings"
)

type (
	MensagemEnviada struct {
		ID             uint64 `json:"id"`
		Mensagem       string `json:"mensagem"`
		RemententeID   uint64 `json:"rementente_id"`
		DestinatarioID uint64 `json:"destinatario_id"`
		CriadoEm       string `json:"criadoem"`
		EnviadoEm      string `json:"enviadoem,omitempty"`
	}
)

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
