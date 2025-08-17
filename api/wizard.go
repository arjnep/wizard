package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Answer struct {
	Text      string
	Archetype string
}

type Question struct {
	Prompt  string
	Answers []Answer
}

type QuestionBank struct {
	Archetypes map[string]string `json:"archetypes"`
	Questions  []Question        `json:"questions"`
}

var qbank QuestionBank

func WizardWave(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Wizard</title>
	</head>
	<body>
		<h1>Hello! Wizard is on the way.</h1>
		{{range $arch, $info := .Archetypes}}
			<ul>
				<li>{{ $arch }}: {{ $info }}</li>
			</ul>
		{{end}}

		{{range $idx, $q := .Questions}}
			<ul>
				<li>{{ $q.Prompt }}</li>
				{{range $a := $q.Answers}}
					<li>{{ $a.Text }}</li>
				{{end}}
			</ul>
		{{end}}
	</body>
	</html>
	`

	dir, _ := os.Getwd()
	fmt.Println("Current dir", dir)

	jsonData, err := os.ReadFile(filepath.Join(dir, "data", "questions.json"))
	if err != nil {
		http.Error(w, "Error reading questions.json"+err.Error(), http.StatusInternalServerError)

		entries, err := os.ReadDir("./")
		if err != nil {
			http.Error(w, "Error reading dir"+err.Error(), http.StatusInternalServerError)
		}

		for _, e := range entries {
			fmt.Println(e.Name())
		}

		return
	}
	if err := json.Unmarshal(jsonData, &qbank); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	templ, err := template.New("wizard").Parse(htmlContent)
	if err != nil {
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}
	if err := templ.Execute(w, qbank); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
		return
	}
}
