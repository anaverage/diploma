package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"statusPage/internal/entities"
	"statusPage/internal/resultData"
)

func Build(router *chi.Mux, store *resultData.ResultDataStorage) {
	router.Use(middleware.Recoverer)

	controller := NewController(store)

	router.Get("/", controller.GetData)

}

type Controller struct {
	storage *resultData.ResultDataStorage
}

func NewController(storage *resultData.ResultDataStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (c *Controller) GetData(w http.ResponseWriter, r *http.Request) {
	var result entities.ResultT
	resultSetT, err := c.storage.GetResultData()
	if err != nil {
		result.Status = false
		result.Error = "Error on collect data"
	} else {
		checkFull := c.storage.IsFull()
		switch checkFull {
		case true:
			result.Status = true
			result.Data = resultSetT
		case false:
			result.Status = false
			result.Error = "Error on collect data"
		}
	}

	// Добавляем вывод сообщения об ошибке
	if result.Error != "" {
		log.Printf("Ошибка: %s", result.Error)
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Printf("Ошибка преобразования ResultT в json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(res)
}
