package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"math"
)

func main() {
	handleReq()
}

func handleReq() {
	r := mux.NewRouter()
	r.HandleFunc("/sensors/cars", getCarSensorData(sendCarData)).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getCarSensorData(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.Open("data/automaticdata.txt")

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer content.Close()
		byteValue, _ := ioutil.ReadAll(content)
		var objs interface{}
		_ = json.Unmarshal([]byte(byteValue), &objs)


		for _, obj := range objs.([]interface{}) {
			for keyProp, valProp := range obj.(map[string]interface{}) {
				if keyProp == "startedat" {
					sec, dec := math.Modf(valProp.(float64))
					obj.(map[string]interface{})[keyProp] = time.Unix(int64(sec), int64(dec*(1e9)))
				} else if keyProp == "endedat" {
					sec, dec := math.Modf(valProp.(float64))
					obj.(map[string]interface{})[keyProp] = time.Unix(int64(sec), int64(dec*(1e9)))
				}
			}
		}

		context.Set(r, "objs", objs)

		h.ServeHTTP(w, r)
	})
}


func sendCarData(w http.ResponseWriter, r *http.Request) {
	allData := context.GetAll(r)
	objs := allData["objs"]
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(objs)
}