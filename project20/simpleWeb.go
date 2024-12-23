package simpleWeb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logName string = "project20/test.log"
var storageName string = "project20/storage.json"

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "server ok"})

}

func logInit() (*os.File, error) {
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	return file, nil
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Method, c.Request.URL.Path)
		// c.JSON(http.StatusOK, gin.H{"status": "server ok"})

	}
}

func storageInit() (*int64, *map[int]Note, *os.File, error) {
	storage, err := os.Open(storageName)
	if err != nil {
		log.Println(err)
		return nil, nil, nil, err
	}

	var maxId int64 = 0
	noteList := make(map[int]Note)

	if fileinfo, _ := os.Stat(storageName); fileinfo.Size() == 0 {
		return &maxId, &noteList, storage, nil
	}

	err = json.NewDecoder(storage).Decode(&noteList)
	if err != nil {
		log.Println(err)
		return nil, nil, nil, err
	}

	for _, i := range noteList {
		if i.Id > maxId {
			maxId = i.Id
		}
	}

	return &maxId, &noteList, storage, nil
}

func storageSave(all *map[int]Note) error {
	file, err := os.OpenFile(storageName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	if len(*all) == 0 {
		log.Println("Попытка сохранения пустой структуры.")
		return errors.New("Попытка сохранения пустой структуры.")
	}
	err = json.NewEncoder(file).Encode(all)
	return err
}

type validator interface {
	isValid() bool
}

func validate(v validator) bool {
	return v.isValid()
}

func (note Note) isValid() bool {
	if note.Description == "" || note.Name == "" {
		return false
	}
	return true
}

func addNote(maxId *int64, noteList *map[int]Note) gin.HandlerFunc {
	return func(c *gin.Context) {
		var note Note
		err := c.ShouldBindJSON(&note)
		log.Println(note)
		if err != nil {
			log.Println("Corrupted json passed.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Corrupted json passed"})
			return
		}

		if !validate(note) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "получены некорректные данные"})
			return
		}
		note.Id = *maxId
		(*noteList)[int(*maxId)] = note
		*maxId++
		storageSave(noteList)
		c.JSON(http.StatusOK, gin.H{"note": note})
	}
}

func showNote(maxId *int64, noteList *map[int]Note) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("Id")

		paramVal, err := strconv.Atoi(param)
		if err != nil || paramVal > int(*maxId) {
			log.Println("Передан некорректный параметр.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан некорректный параметр"})

			return
		}

		if (*noteList)[paramVal].Id == 0 {
			log.Println("Элемент не найден.")
			c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"note": (*noteList)[paramVal]})

	}
}

func showAll(noteList *map[int]Note) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len((*noteList)) == 0 {
			log.Println("Нет сохранённых заметок.")
			c.JSON(http.StatusNoContent, gin.H{"error": "Нет сохранённых заметок."})

			return
		}

		c.JSON(http.StatusOK, gin.H{"note": (*noteList)})

	}
}

func deleteNote(maxId *int64, noteList *map[int]Note) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")

		paramVal, err := strconv.Atoi(param)
		if err != nil || paramVal > int(*maxId) {
			log.Println("Передан некорректный параметр.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан некорректный параметр"})

			return
		}

		if (*noteList)[paramVal].Id == 0 {
			log.Println("Элемент не найден.")
			c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})

			return
		}

		var note Note = (*noteList)[paramVal]
		delete((*noteList), paramVal)
		storageSave(noteList)
		c.JSON(http.StatusOK, gin.H{"deleted note": note})

	}
}

func updateNote(maxId *int64, noteList *map[int]Note) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")

		paramVal, err := strconv.Atoi(param)
		if err != nil || paramVal > int(*maxId) {
			log.Println("Передан некорректный параметр.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Передан некорректный параметр"})

			return
		}

		if (*noteList)[paramVal].Id == 0 {
			log.Println("Элемент не найден.")
			c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})

			return
		}
		var note Note
		err = c.ShouldBindJSON(&note)
		if err != nil || !validate(note) {
			log.Println("Corrupted json passed.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Corrupted json passed"})

			return
		}

		note.Id = int64(paramVal)

		(*noteList)[paramVal] = note
		storageSave(noteList)

		c.JSON(http.StatusOK, gin.H{"updated note": note})

	}
}

func WebServer() {
	Log, err := logInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Log.Close()

	maxId, noteList, storage, err := storageInit()
	if err != nil {
		log.Println(err)
		return
	}
	defer storage.Close()
	if *maxId == 0 {
		*maxId = 1
	}
	router := gin.Default()

	router.GET("/ping", logMiddleware(), ping)
	router.POST("/Note", logMiddleware(), addNote(maxId, noteList))
	router.GET("/Note/:id", logMiddleware(), showNote(maxId, noteList))
	router.DELETE("/Note/:id", logMiddleware(), deleteNote(maxId, noteList))
	router.PUT("/Note/:id", logMiddleware(), updateNote(maxId, noteList))
	router.GET("/Notes", logMiddleware(), showAll(noteList))

	router.Run("localhost:8081")
}
