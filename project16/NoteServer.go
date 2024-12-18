package NoteServer

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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
		(*noteLst)[note.Id] = note
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

func deleteNote(noteLst *map[int]Note, mId *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var note Note
		if r.Method != http.MethodDelete {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		params := r.URL.Query()
		var buf int
		var err error

		buf, err = strconv.Atoi(params.Get("Id"))

		if err != nil {
			http.Error(w, "Id is missing", http.StatusBadRequest)
			return
		}

		if (*noteLst)[buf].Id == 0 || buf > *mId {
			http.Error(w, "Id not found", http.StatusNotFound)
			return
		}

		delete(*noteLst, buf)
		storageSave(noteLst)

		w.WriteHeader(http.StatusGone)
	}
}

func updateNote(noteLst *map[int]Note) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		params := r.URL.Query()

		buf, err := strconv.Atoi(params.Get("Id"))

		if err != nil {
			http.Error(w, "Param is missing", http.StatusBadRequest)
			return
		}
		var note Note

		err = json.NewDecoder(r.Body).Decode(&note)

		if err != nil {
			http.Error(w, "Json is corrupted", http.StatusBadRequest)
			return
		}
		if !isValid(&note) {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}
		if (*noteLst)[buf].Id == 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		note.Id = buf

		(*noteLst)[buf] = note
		storageSave(noteLst)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	}
}

func showNote(noteLst *map[int]Note) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		params := r.URL.Query()

		buf, err := strconv.Atoi(params.Get("Id"))

		if err != nil {
			http.Error(w, "Param is missing", http.StatusBadRequest)
			return
		}

		if (*noteLst)[buf].Id == 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode((*noteLst)[buf])
	}
}

func allNotes(noteLst *map[int]Note) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(noteLst)
	}
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
	h.HandleFunc("/Notes/delete", deleteNote(&noteLst, &noteCurId))
	h.HandleFunc("/Notes/update", updateNote(&noteLst))
	h.HandleFunc("/Notes/show", showNote(&noteLst))
	h.HandleFunc("/Notes", allNotes(&noteLst))

	http.ListenAndServe("localhost:8080", h)
}
