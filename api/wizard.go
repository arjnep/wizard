package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
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

	jsonData, err := os.ReadFile("api/questions.json")
	if err != nil {
		http.Error(w, "Error reading questions.json", http.StatusInternalServerError)
		dir, _ := os.Getwd()
		fmt.Println("Error :", err, "Current dir", dir)
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
