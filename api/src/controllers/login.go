package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respostas.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var usuario modelos.Usuario
	if err := json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, err := repositorio.BuscarPorEmail(usuario.Email)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err := seguranca.VerificaSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	usuarioID := strconv.FormatUint(usuarioSalvoNoBanco.ID, 10)
	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioID, Token: token})
}
