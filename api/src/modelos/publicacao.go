package modelos

import (
	"errors"
	"strings"
	"time"
)

type Publicacao struct {
	ID        uint64    `json:"id,omitempty"`
	Titulo    string    `json:"titulo,omitempty"`
	Conteudo  string    `json:"conteudo,omitempty"`
	AutorID   uint64    `json:"autorid,omitempty"`
	AutorNick string    `json:"autornick,omitempty"`
	Curtidas  uint64    `json:"curtidas"`
	CriadoEm  time.Time `json:"criadoem,omitempty"`
}

func (p *Publicacao) Preparar() error {
	if err := p.validar(); err != nil {
		return err
	}

	if err := p.formatar(); err != nil {
		return err
	}

	return nil
}

func (p *Publicacao) validar() error {
	if p.Titulo == "" {
		return errors.New("o título é obrigatório e não pode estar em branco")
	}

	if p.Conteudo == "" {
		return errors.New("o conteúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (p *Publicacao) formatar() error {
	p.Titulo = strings.TrimSpace(p.Titulo)
	p.Conteudo = strings.TrimSpace(p.Conteudo)

	return nil
}
