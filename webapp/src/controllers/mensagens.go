package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/requisicoes"
	"webapp/src/respostas"

	"github.com/gorilla/mux"
)

func EnviarMensagem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mensagem, err := json.Marshal(map[string]string{
		"mensagem": r.FormValue("mensagem"),
	})
	if err != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%d/mensagens", config.ApiUrl, usuarioID)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, bytes.NewBuffer(mensagem))
	if err != nil {
		fmt.Println("FazerRequisicaoComAutenticacao")
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Println("StatusCode")
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	respostas.JSON(w, response.StatusCode, nil)
}
