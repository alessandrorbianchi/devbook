package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasMensagens = []Rota{
	{
		Uri:                "/mesangens",
		Metodo:             http.MethodPost,
		Funcao:             controllers.EnviarMensagem,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/mensagens",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMensagens,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/usuarios/{usuarioId}/mensagens",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMensagensPorUsuario,
		RequerAutenticacao: true,
	},
	// {
	// 	Uri:                "/mensagens/{mensagemId}",
	// 	Metodo:             http.MethodGet,
	// 	Funcao:             controllers.BuscarMensagem,
	// 	RequerAutenticacao: true,
	// },
	// {
	// 	Uri:                "/publicacoes/{publicacaoId}",
	// 	Metodo:             http.MethodDelete,
	// 	Funcao:             controllers.DeletarPublicacao,
	// 	RequerAutenticacao: true,
	// },

}
