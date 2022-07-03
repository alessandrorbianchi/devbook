package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

type Usuarios struct {
	db *sql.DB
}

func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

func (u Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	statement, err := u.db.Prepare("insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	resultado, err := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		return 0, err
	}

	ultimoIDInserido, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoIDInserido), nil
}

func (u Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)

	linhas, err := u.db.Query(
		"SELECT id, nome, nick, email, criadoem FROM usuarios WHERE nome LIKE ? OR nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (u Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linhas, err := u.db.Query(
		"SELECT id, nome, nick, email, criadoem FROM usuarios WHERE id = ?",
		ID,
	)
	if err != nil {
		return modelos.Usuario{}, err
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return modelos.Usuario{}, err
		}
	}

	return usuario, nil
}

func (u Usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	statement, err := u.db.Prepare(
		"UPDATE usuarios SET nome = ?, nick = ?, email = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); err != nil {
		return err
	}

	return nil
}

func (u Usuarios) Deletar(ID uint64) error {
	statement, err := u.db.Prepare("DELETE FROM usuarios WHERE ID = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}

	return nil
}

func (u Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, err := u.db.Query("SELECT id, senha FROM usuarios WHERE email = ?", email)
	if err != nil {
		return modelos.Usuario{}, err
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if err := linha.Scan(&usuario.ID, &usuario.Senha); err != nil {
			return modelos.Usuario{}, err
		}
	}

	return usuario, nil
}

func (u Usuarios) Seguir(usuarioID, seguidorID uint64) error {
	statement, err := u.db.Prepare(
		"INSERT IGNORE INTO seguidores (usuario_id, seguidor_id) VALUES (?, ?)",
	)

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioID, seguidorID); err != nil {
		return err
	}

	return nil
}

func (u Usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	statement, err := u.db.Prepare("DELETE FROM seguidores WHERE usuario_id = ? and seguidor_id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioID, seguidorID); err != nil {
		return err
	}

	return nil
}

func (u Usuarios) BuscarSeguidores(usuarioID uint64) ([]modelos.Usuario, error) {
	linhas, err := u.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem 
			FROM usuarios u INNER JOIN seguidores s ON u.id = s.seguidor_id
			WHERE s.usuario_id = ?
		`, usuarioID)

	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (u Usuarios) BuscarSeguindo(usuarioID uint64) ([]modelos.Usuario, error) {
	linhas, err := u.db.Query(`
		SELECT u.id, u.nome, u.nick, u.email, u.criadoem 
			FROM usuarios u INNER JOIN seguidores s ON u.id = s.usuario_id
			WHERE s.seguidor_id = ?
		`, usuarioID)

	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if err = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (u Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, err := u.db.Query("SELECT senha FROM usuarios WHERE id = ?", usuarioID)
	if err != nil {
		return "", err
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if err := linha.Scan(&usuario.Senha); err != nil {
			return "", err
		}
	}

	return usuario.Senha, nil
}

func (u Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, err := u.db.Prepare(
		"UPDATE usuarios SET senha = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(senha, usuarioID); err != nil {
		return err
	}

	return nil
}
