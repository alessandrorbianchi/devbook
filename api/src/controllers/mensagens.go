package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func EnviarMensagem(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	corpoRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respostas.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var mensagem modelos.MensagemEnviada
	if err = json.Unmarshal(corpoRequest, &mensagem); err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	mensagem.RemententeID = usuarioID
	if err := mensagem.Preparar(); err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeMensagens(db)
	mensagem.ID, err = repositorio.Enviar(mensagem)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusCreated, mensagem)
}

func BuscarMensagens(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeMensagens(db)
	mensagens, err := repositorio.Buscar(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, mensagens)
}

func BuscarMensagensPorUsuario(w http.ResponseWriter, r *http.Request) {
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

	repositorio := repositorios.NovoRepositorioDeMensagens(db)
	mensagens, err := repositorio.BuscarPorUsuario(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, mensagens)
}

// func BuscarMensagem(w http.ResponseWriter, r *http.Request) {
// 	parametros := mux.Vars(r)

// 	mensagemID, err := strconv.ParseUint(parametros["mensagemID"], 10, 64)
// 	if err != nil {
// 		respostas.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := banco.Conectar()
// 	if err != nil {
// 		respostas.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repositorio := repositorios.NovoRepositorioDeMensagens(db)
// 	publicacao, err := repositorio.BuscarPorID(mensagemID)
// 	if err != nil {
// 		respostas.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	respostas.JSON(w, http.StatusOK, publicacao)
// }

// func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
// 	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
// 	if err != nil {
// 		respostas.Error(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	parametros := mux.Vars(r)
// 	publicacaoID, err := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
// 	if err != nil {
// 		respostas.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := banco.Conectar()
// 	if err != nil {
// 		respostas.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
// 	publicacaoSalvaNoBanco, err := repositorio.BuscarPorID(publicacaoID)
// 	if err != nil {
// 		respostas.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	if publicacaoSalvaNoBanco.AutorID != usuarioID {
// 		respostas.Error(w, http.StatusForbidden, fmt.Errorf("não é possível deletar uma publicação que não é sua"))
// 		return
// 	}

// 	if err = repositorio.Deletar(publicacaoID); err != nil {
// 		respostas.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	respostas.JSON(w, http.StatusNoContent, nil)
// }
