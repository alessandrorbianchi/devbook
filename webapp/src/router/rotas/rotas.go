package rotas

import (
	"net/http"
	"webapp/src/middlewares"

	"github.com/gorilla/mux"
)

type Rota struct {
	Uri                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasLogin
	rotas = append(rotas, rotasUsuarios...)
	rotas = append(rotas, rotaPaginaPrincipal)
	rotas = append(rotas, rotasPublicacoes...)
	rotas = append(rotas, rotaLogout)
	rotas = append(rotas, rotasMensagens...)

	for _, rota := range rotas {
		if rota.RequerAutenticacao {
			r.HandleFunc(rota.Uri,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.Uri,
				middlewares.Logger(rota.Funcao),
			).Methods(rota.Metodo)
		}
	}

	fileServer := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	return r
}
