package modelos

import "time"

type Publicacao struct {
	ID        uint64    `json:"id,omitempty"`
	Titulo    string    `json:"titulo,omitempty"`
	Conteudo  string    `json:"conteudo,omitempty"`
	AutorID   uint64    `json:"autorid,omitempty"`
	AutorNick string    `json:"autornick,omitempty"`
	Curtidas  uint64    `json:"curtidas"`
	CriadoEm  time.Time `json:"criadoem,omitempty"`
}
