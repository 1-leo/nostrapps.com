package main

import (
	"fmt"
	"net/http"
	"strings"
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

type AppDefinition struct {
	Name        string   `yaml:"name"`
	Npub        string   `yaml:"npub"`
	Categories  []string `yaml:"categories"`
	Description string   `yaml:"description"`
	Gallery     []string `yaml:"gallery"`
	Features    []string `yaml:"features"`
	Platforms   []string `yaml:"platforms"`
	Source      string   `yaml:"source"`
	Thumb       string   `yaml:"thumb"`
	URL         string   `yaml:"url"`
}

func (a AppDefinition) Fulltext() string {
	return strings.Join([]string{
		a.Name,
		a.Description,
		a.URL,
		a.Source,
	}, " ")
}
