package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {
		respostas.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var usuario modelos.Usuario
	if err = json.Unmarshal(corpoRequest, &usuario); err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = usuario.Preparar("cadastro"); err != nil {
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
	usuario.ID, err = repositorio.Criar(usuario)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, err := repositorio.Buscar(nomeOuNick)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
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
	usuario, err := repositorio.BuscarPorID(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Error(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respostas.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var usuario modelos.Usuario
	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := usuario.Preparar("edicao"); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if err := repositorio.Atualizar(usuarioID, usuario); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Error(w, http.StatusForbidden, errors.New("não é possível deletar um usuário que não seja o seu"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if err = repositorio.Deletar(usuarioID); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	if seguidorID == usuarioID {
		respostas.Error(w, http.StatusForbidden, errors.New("não é ppossível seguir você mesmo"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if err = repositorio.Seguir(usuarioID, seguidorID); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	if seguidorID == usuarioID {
		respostas.Error(w, http.StatusForbidden, errors.New("não é ppossível parar de seguir você mesmo"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if err = repositorio.PararDeSeguir(usuarioID, seguidorID); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
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
	seguidores, err := repositorio.BuscarSeguidores(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)
}

func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
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
	seguindo, err := repositorio.BuscarSeguindo(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, seguindo)
}

func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Error(w, http.StatusForbidden, errors.New("não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respostas.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var senha modelos.Senha
	if err = json.Unmarshal(corpoRequisicao, &senha); err != nil {
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
	senhaSalvaNoBanco, err := repositorio.BuscarSenha(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = seguranca.VerificaSenha(senhaSalvaNoBanco, senha.Atual); err != nil {
		respostas.Error(w, http.StatusUnauthorized, errors.New("a senha atual não condiz com a que está salva no banco"))
		return
	}

	senhaComHash, err := seguranca.Hash(senha.Nova)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
