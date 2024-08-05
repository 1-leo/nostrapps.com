package main

import (
	"fmt"
	"net/http"
)

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	home(categories, platforms, apps).Render(r.Context(), w)
}

func handleAppPage(w http.ResponseWriter, r *http.Request) {
	appId := r.PathValue("app")
	def, ok := apps[appId]
	if !ok {
		http.Error(w, fmt.Sprintf("app %s not known", appId), 404)
		return
	}

	app(platforms, def).Render(r.Context(), w)
}
