package handler

import (
	"net/http"
	"text/template"

	"github.com/MihaiBlebea/go-wallet/domain"
)

func WebappEndpoint(tmpl *template.Template, wallet domain.Wallet) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := r.FormValue("start")
		end := r.FormValue("end")

		data := struct {
			Start string
			End   string
		}{
			Start: start,
			End:   end,
		}

		w.WriteHeader(200)
		tmpl.ExecuteTemplate(w, "report", data)
	})
}
