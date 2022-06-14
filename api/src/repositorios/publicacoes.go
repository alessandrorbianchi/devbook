package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

func (p Publicacoes) Criar(publicacao modelos.Publicacao) (uint64, error) {
	statement, err := p.db.Prepare("INSERT INTO publicacoes (titulo, conteudo, autor_id) values (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	resultado, err := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AutorID)
	if err != nil {
		return 0, err
	}

	ultimoIDInserido, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoIDInserido), nil
}

func (p Publicacoes) Buscar(usuarioID uint64) ([]modelos.Publicacao, error) {
	linhas, err := p.db.Query(
		`SELECT DISTINCT p.id, p.titulo, p.conteudo, p.autor_id, p.curtidas, p.criadoem, u.nick 
			FROM publicacoes p 
			INNER JOIN usuarios u ON u.id = p.autor_id
			INNER JOIN seguidores s ON p.autor_id = s.usuario_id
		WHERE u.id = ? or s.seguidor_id = ? 
		order by 1 desc`,
		usuarioID, usuarioID,
	)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if err = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (p Publicacoes) BuscarPorID(publicacaoID uint64) (modelos.Publicacao, error) {
	linha, err := p.db.Query(
		`SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.curtidas, p.criadoem, u.nick 
			FROM publicacoes p INNER JOIN usuarios u ON 
			u.id = p.autor_id WHERE p.id = ?`,
		publicacaoID,
	)
	if err != nil {
		return modelos.Publicacao{}, err
	}
	defer linha.Close()

	var publicacao modelos.Publicacao

	if linha.Next() {
		if err = linha.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); err != nil {
			return modelos.Publicacao{}, err
		}
	}

	return publicacao, nil
}

func (p Publicacoes) Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error {
	statement, err := p.db.Prepare(
		"UPDATE publicacoes SET titulo = ?, conteudo = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoID); err != nil {
		return err
	}

	return nil
}

func (p Publicacoes) Deletar(publicacaoID uint64) error {
	statement, err := p.db.Prepare("DELETE FROM publicacoes WHERE id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoID); err != nil {
		return err
	}

	return nil
}

func (p Publicacoes) BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error) {
	linha, err := p.db.Query(
		`SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.curtidas, p.criadoem, u.nick 
			FROM publicacoes p 
			INNER JOIN usuarios u ON u.id = p.autor_id 
			WHERE p.autor_id = ?`,
		usuarioID,
	)
	if err != nil {
		return nil, err
	}
	defer linha.Close()

	var publicacoes []modelos.Publicacao

	for linha.Next() {
		var publicacao modelos.Publicacao

		if err = linha.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (p Publicacoes) Curtir(publicacaoID uint64) error {
	statement, err := p.db.Prepare(
		"UPDATE publicacoes SET curtidas = curtidas + 1 WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoID); err != nil {
		return err
	}

	return nil
}

func (p Publicacoes) Descurtir(publicacaoID uint64) error {
	statement, err := p.db.Prepare(
		`UPDATE publicacoes SET curtidas = 
			CASE 
				WHEN curtidas > 0 THEN curtidas - 1
				ELSE 0 
			END 
			WHERE id = ?`,
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoID); err != nil {
		return err
	}

	return nil
}
