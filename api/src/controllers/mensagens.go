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

	"github.com/gorilla/mux"
)

func EnviarMensagem(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	destinatarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
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
	mensagem.DestinatarioID = destinatarioID
	codigoDeSeguranca, err := seguranca.GerarCodigoDeSeguranca(usuarioID, destinatarioID)
	if err != nil {
		respostas.Error(w, http.StatusBadRequest, err)
		return
	}
	mensagem.CodigoSeguranca = codigoDeSeguranca
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

func BuscarMensagensAgrupadasPorUsuario(w http.ResponseWriter, r *http.Request) {
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
	mensagens, err := repositorio.BuscarAgrupado(usuarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, mensagens)
}

func BuscarMensagensPorUsuario(w http.ResponseWriter, r *http.Request) {
	remententeID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Error(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	destinatarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
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
	mensagens, err := repositorio.BuscarPorUsuario(remententeID, destinatarioID)
	if err != nil {
		respostas.Error(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, mensagens)
}
