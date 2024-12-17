package NoteServer

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const logFile = "project16/test.log"
const storageFile = "project16/storage.json"

type Note struct {
	Id          int    `json:"Id"`
	StatusRate  int    `json:"StatusRate"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func isValid(note *Note) bool {
	if (*note).Description == "" || (*note).Name == "" || (*note).StatusRate < 0 {
		return false
	}
	return true
}

func addNote(noteLst *map[int]Note, mId *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note Note
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&note)

		if err != nil || !isValid(&note) {
			http.Error(w, "Passed invalid json", http.StatusBadRequest)
			return
		}
		note.Id = *mId
		*mId++
		(*noteLst)[*mId] = note
		err = storageSave(noteLst)
		if err != nil {
			http.Error(w, "Passed invalid json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func deleteNote(w http.ResponseWriter, r *http.Request) {

}

func updateNote(w http.ResponseWriter, r *http.Request) {

}

func showNote(w http.ResponseWriter, r *http.Request) {

}

func allNotes(w http.ResponseWriter, r *http.Request) {

}

func logPrepare() (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.SetOutput(file)
	return file, err
}

func storageInit(noteLst *map[int]Note, id *int) error {
	if fileInfo, _ := os.Stat(storageFile); fileInfo.Size() == 0 {
		log.Println("Хранилище пусто.")
		return nil
	}

	file, err := os.Open(storageFile)

	if err != nil {
		log.Println(err)
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&noteLst)

	if err != nil {
		log.Println(err)
		return err
	}
	for _, i := range *noteLst {
		*id = i.Id + 1
	}

	return nil
}

func storageSave(noteLst *map[int]Note) error {
	file, err := os.OpenFile(storageFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	encoder := json.NewEncoder(file)

	encoder.Encode(noteLst)

	return nil
}

func ToDoServer() {

	file, err := logPrepare()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	// В других проектах сделал через глобальные переменные, сейчас переделаю без них.
	var noteLst = make(map[int]Note)
	var noteCurId int = 1

	// Инициализация хранилища.
	err = storageInit(&noteLst, &noteCurId)

	if err != nil {
		log.Println(err)
		return
	}

	h := http.NewServeMux()
	h.HandleFunc("/Notes/add", addNote(&noteLst, &noteCurId))
	h.HandleFunc("/Notes/delete", deleteNote)
	h.HandleFunc("/Notes/update", updateNote)
	h.HandleFunc("/Notes/show", showNote)
	h.HandleFunc("/Notes", allNotes)

	http.ListenAndServe("localhost:8080", h)
}
