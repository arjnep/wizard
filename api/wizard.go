package handler

import (
	"fmt"
	"net/http"
)

func WizardWave(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Wizard</title>
	</head>
	<body>
		<h1>Hello! Wizard is on the way.</h1>
	</body>
	</html>
	`
	fmt.Fprint(w, htmlContent)
}
