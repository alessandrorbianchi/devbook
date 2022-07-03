package rotas

import (
	"net/http"
	"webapp/src/controllers"
)

var rotasMensagens = []Rota{
	{
		Uri:                "/usuarios/{usuarioId}/mensagens",
		Metodo:             http.MethodPost,
		Funcao:             controllers.EnviarMensagem,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/mensagem-usuario",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaDeMensagensDoUsuario,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/mensagem-seguidor/{usuarioId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaDeMensagensDoSeguidor,
		RequerAutenticacao: true,
	},
}
