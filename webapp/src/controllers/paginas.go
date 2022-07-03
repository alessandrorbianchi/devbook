package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/modelos"
	"webapp/src/requisicoes"
	"webapp/src/respostas"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func CarregarTelaDeLogin(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)

	if cookie["token"] != "" {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	utils.ExecutarTemplate(w, "login.html", nil)
}

func CarregarPaginaDeCadastroDeUsuario(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "cadastro.html", nil)
}

func CarregarPaginaPrincipal(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/publicacoes", config.ApiUrl)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var publicacoes []modelos.Publicacao
	if err = json.NewDecoder(response.Body).Decode(&publicacoes); err != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	utils.ExecutarTemplate(w, "home.html", struct {
		Publicacoes []modelos.Publicacao
		UsuarioID   uint64
	}{
		Publicacoes: publicacoes,
		UsuarioID:   usuarioID,
	})
}

func CarregarPaginaDeAtualizacaoDePublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, err := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if err != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d", config.ApiUrl, publicacaoID)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var publicacao modelos.Publicacao
	if err = json.NewDecoder(response.Body).Decode(&publicacao); err != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "atualizar-publicacao.html", publicacao)
}

func CarregarPaginaDeUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	url := fmt.Sprintf("%s/usuarios?usuario=%s", config.ApiUrl, nomeOuNick)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var usuarios []modelos.Usuario
	if err := json.NewDecoder(response.Body).Decode(&usuarios); err != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "usuarios.html", usuarios)
}

func CarregarPerfildDoUsuarios(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioLogadoID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	if usuarioID == usuarioLogadoID {
		http.Redirect(w, r, "/perfil", http.StatusFound)
		return
	}

	usuario, err := modelos.BuscarUsuarioCompleto(usuarioID, r)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "usuario.html", struct {
		Usuario         modelos.Usuario
		UsuarioLogadoID uint64
	}{
		Usuario:         usuario,
		UsuarioLogadoID: usuarioLogadoID,
	})
}

func CarregarPerfilDoUsuarioLogado(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	usuario, err := modelos.BuscarUsuarioCompleto(usuarioID, r)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "perfil.html", usuario)
}

func CarregarPaginaDeEdicaoDeUsuario(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	canal := make(chan modelos.Usuario)
	go modelos.BuscarDadosDoUsuario(canal, usuarioID, r)
	usuario := <-canal

	if usuario.ID == 0 {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: "Erro ao buscar usuário"})
		return
	}

	utils.ExecutarTemplate(w, "editar-usuario.html", usuario)
}

func CarregarPaginaDeAtualizacaoDeSenha(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "atualizar-senha.html", nil)
}

func CarregarPaginaDeMensagensDoUsuario(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/mensagens", config.ApiUrl)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var mensagens []modelos.Mensagem
	if err = json.NewDecoder(response.Body).Decode(&mensagens); err != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	utils.ExecutarTemplate(w, "mensagem-usuario.html", struct {
		Mensagens []modelos.Mensagem
		UsuarioID uint64
	}{
		Mensagens: mensagens,
		UsuarioID: usuarioID,
	})
}

func CarregarPaginaDeMensagensDoSeguidor(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%d/mensagens", config.ApiUrl, usuarioID)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErrorAPI{Error: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var mensagens []modelos.Mensagem
	if err = json.NewDecoder(response.Body).Decode(&mensagens); err != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErrorAPI{Error: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioLogadoID, _ := strconv.ParseUint(cookie["id"], 10, 64)

	utils.ExecutarTemplate(w, "mensagem-seguidor.html", struct {
		Mensagens       []modelos.Mensagem
		UsuarioID       uint64
		UsuarioLogadoID uint64
	}{
		Mensagens:       mensagens,
		UsuarioID:       usuarioID,
		UsuarioLogadoID: usuarioLogadoID,
	})
}
