package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
)

type Note struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	NoteText string `json:"note_text"`
}

var NoteStorage = []Note{}

func main() {
	e := echo.New()
	e.GET("/hello", GetHello)
	e.POST("/save_note", SaveNote)
	e.GET("/list_all", ListAllNotes)
	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}

func GetHello(c echo.Context) error {
	name := c.QueryParam("name")
	ln := c.QueryParam("ln")
	return c.String(200, "Привет, " + name + " " + ln)
}

func SaveNote(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}

	newNote := Note{}
	err = json.Unmarshal(body, &newNote)
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}
	NoteStorage = append(NoteStorage, newNote)

	fmt.Printf("Введённые данные: \n  имя: %s\n  фамилия: %s\n  заметка: %s\n", newNote.Name, newNote.Surname, newNote.NoteText)
	return c.JSON(200, json.RawMessage(body))
}

func ListAllNotes(c echo.Context) error {
	jsonResp, err := json.Marshal(NoteStorage)
	if err != nil {
		log.Println(err)
		return c.NoContent(500)
	}

	return c.JSON(200, json.RawMessage(jsonResp))
}
