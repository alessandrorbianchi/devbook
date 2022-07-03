package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasMensagens = []Rota{
	{
		Uri:                "/usuarios/{usuarioId}/mensagens",
		Metodo:             http.MethodPost,
		Funcao:             controllers.EnviarMensagem,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/mensagens",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMensagensAgrupadasPorUsuario,
		RequerAutenticacao: true,
	},
	{
		Uri:                "/usuarios/{usuarioId}/mensagens",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMensagensPorUsuario,
		RequerAutenticacao: true,
	},
}
