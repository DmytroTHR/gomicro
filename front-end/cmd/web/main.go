package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const webPort = "8088"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port", webPort)
	err := http.ListenAndServe(":"+webPort, nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	templateSlice := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	templateSlice = append([]string{fmt.Sprintf("templates/%s", t)}, templateSlice...)

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTmpl struct {
		BrokerURL string
	}

	// brokerHost := "10.109.224.181"
	// brokerPort := "8080"
	// dataTmpl.BrokerURL = "http://" + brokerHost + ":" + brokerPort

	dataTmpl.BrokerURL = os.Getenv("BROKER_URL")

	if err := tmpl.Execute(w, dataTmpl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
