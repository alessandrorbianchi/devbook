package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoem,omitempty"`
}

func (u *Usuario) Preparar(etapa string) error {
	if err := u.validar(etapa); err != nil {
		return err
	}

	if err := u.formatar(etapa); err != nil {
		return err
	}

	return nil
}

func (u *Usuario) validar(etapa string) error {
	if u.Nome == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}

	if u.Nick == "" {
		return errors.New("o nick é obrigatório e não pode estar em branco")
	}

	if u.Email == "" {
		return errors.New("o email é obrigatório e não pode estar em branco")
	}

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("o email inserido é inválido")
	}

	if etapa == "cadastro" && u.Senha == "" {
		return errors.New("a senha é obrigatória e não pode estar em branco")
	}

	return nil
}

func (u *Usuario) formatar(etapa string) error {
	u.Nome = strings.TrimSpace(u.Nome)
	u.Nick = strings.TrimSpace(u.Nick)
	u.Email = strings.TrimSpace(u.Email)

	if etapa == "cadastro" {
		senhaComHash, err := seguranca.Hash(u.Senha)
		if err != nil {
			return err
		}

		u.Senha = string(senhaComHash)
	}

	return nil
}
