package thisapp

import (
	"context"
	"time"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
)

type Response struct {
	ServerProcessTime string      `json:"server_process_time"`
	ErrorMessage      []string    `json:"message_error,omitempty"`
	StatusMessage     []string    `json:"message_status,omitempty"`
	Data interface{} `json:"data"`
}

type ResponseJSON func(rw http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn ResponseJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Response Format
	response := Response{}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,OPTIONS,POST")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	//Init Context
	ctx := r.Context()

	//Set a 10 second timeout on responses. TODO: make this come from config
	ctx, cancelFn := context.WithTimeout(ctx, 10 * time.Second)
	defer cancelFn()

	//Add context into HTTP Request
	r = r.WithContext(ctx)

	//Start Timer
	start := time.Now()

	//Do the Function
	data, err := fn(w, r)

	//Record Elapsed Time
	response.ServerProcessTime = time.Since(start).String()

	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		response.Data = data
		if buf, err := json.Marshal(response); err == nil {
			w.Write(buf)
			return
		}
	}

	if err != nil {
		response.ErrorMessage = []string{
			err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	buf, _ := json.Marshal(response)
	w.Write(buf)
	return
}


func (am *AppModule) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This App is Alive"))
}

func (am *AppModule) SuccessAPIHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	dt := map[string]string{
		"A": "a",
		"B": "b",
	}
	return dt, nil
}

func (am *AppModule) ErrorAPIHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, fmt.Errorf("Failed as Default")
}

func (am *AppModule) DatabaseAPIHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	products, err := am.ProductDB.GetAll(ctx, "product_id", "desc")
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("Gagal get data dari database")
	}

	return products, nil
}

func (am *AppModule) InitHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/thisapp", am.DefaultHandler)
	mux.Handle("/api/thisapp/success", ResponseJSON(am.SuccessAPIHandler))
	mux.Handle("/api/thisapp/error", ResponseJSON(am.ErrorAPIHandler))
	mux.Handle("/api/thisapp/db", ResponseJSON(am.DatabaseAPIHandler))
}
