package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Note struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	NoteText string `json:"note_text"`
}

type Notification struct {
	Type string
	Text string
}

var reader = bufio.NewReader(os.Stdin)
var httpClient = &http.Client{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		fmt.Println("Parsing file error ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	notification := makeNotification(r.URL.Query().Get("message"))
	err = t.ExecuteTemplate(w, "index.html", &notification)
	if err != nil {
		fmt.Println("Unable to execute template ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func save_note(w http.ResponseWriter, r *http.Request) {
	NewNote := Note{}
	NewNote.Name = r.FormValue("first_name")
	NewNote.Surname = r.FormValue("last_name")
	NewNote.NoteText = r.FormValue("note_text")

	json, err := json.Marshal(NewNote)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8080/save_note", "application/json", bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/?message=error", http.StatusFound)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		http.Redirect(w, r, "/?message=error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/?message=created", http.StatusFound)
	return
}

func list_all(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:8080/list_all")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	var noteList []Note
	if json.Unmarshal(body, &noteList) != nil {
		fmt.Println("Error in parsing json ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t := template.New("list.html").Funcs(template.FuncMap{"inc": func(i int) int { return i + 1 }})
	t, err = t.ParseFiles("static/list.html")
	if err != nil {
		fmt.Println("Parsing file error ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "list.html", noteList)
	if err != nil {
		fmt.Println("Unable to execute template ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func delete_note(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	noteId := mux.Vars(r)["id"]
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/delete_note/%s", noteId))
	if err != nil {
		fmt.Println("error in request: ", err)
		http.Redirect(w, r, "/?message=error", http.StatusFound)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		http.Redirect(w, r, "/?message=error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/list_all", http.StatusFound)
	return
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("./static/stylesheets"))))
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/save_note", save_note)
	router.HandleFunc("/list_all", list_all)
	router.HandleFunc("/delete_note/{id:[0-9]+}", delete_note)
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func makeNotification(message string) *Notification {
	newNotification := Notification{}
	if message == "" {
		return &newNotification
	}

	switch message {
	case "created":
		newNotification.Type = "success"
		newNotification.Text = "Your note was successfully created!"
	case "deleted":
		newNotification.Type = "success"
		newNotification.Text = "Your note was successfully deleted!"
	case "error":
		newNotification.Type = "danger"
		newNotification.Text = "Something went wrong. Try again."
	}
	return &newNotification
}
