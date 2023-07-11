package main

import (
	"encoding/json"
	_ "encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	_ "os"
)

type Message struct {
	Title string
	Body  string
}

var Messages []Message

func saveMessagesToFile() error {
	messagesFile, err := os.Create("messages.json")
	if err != nil {
		return err
	}
	defer func(messagesFile *os.File) {
		err := messagesFile.Close()
		if err != nil {

		}
	}(messagesFile)

	return json.NewEncoder(messagesFile).Encode(Messages)
}

func loadMessagesFromFile() error {
	messagesFile, err := os.Open("messages.json")
	if err != nil {
		if os.IsNotExist(err) { // Caso o arquivo não exista, retorna nil indicando que não é um erro.
			return nil
		}
		// Se o erro for qualquer outro, retorna o erro.
		return err
	}
	defer func(messagesFile *os.File) {
		err := messagesFile.Close()
		if err != nil {

		}
	}(messagesFile)

	// Se o arquivo estiver vazio, não faça a leitura.
	fileInfo, err := messagesFile.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return nil
	}

	// Se o arquivo não estiver vazio, faça a leitura.
	return json.NewDecoder(messagesFile).Decode(&Messages)
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t, err := template.ParseFiles("home.html")
		if err != nil {
			log.Print("Erro ao fazer parse do template: ", err)
		}

		err = t.Execute(w, Messages)
		if err != nil {
			return
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return
		}

		newMessage := Message{
			Title: r.FormValue("title"),
			Body:  r.FormValue("body"),
		}

		Messages = append(Messages, newMessage)
		err = saveMessagesToFile()
		if err != nil {
			log.Print("Erro ao salvar as mensagens: ", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	err := loadMessagesFromFile()
	if err != nil {
		log.Print("Erro ao carregar as mensagens: ", err)
	}

	http.HandleFunc("/", HomePage)
	http.Handle("/styles.css", http.FileServer(http.Dir(".")))
	log.Print("Iniciando o servidor na porta 8080")
	log.Print("Acesse http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
