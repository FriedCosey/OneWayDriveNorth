package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	handleReq()
}

func handleReq() {
	r := mux.NewRouter()
	r.HandleFunc("/sensors/cars", getCarSensorData(sendData)).Methods("GET", "OPTIONS")
	r.HandleFunc("/sensors/microwave", getMicroWaveSensorData(sendData)).Methods("GET", "OPTIONS")
	r.HandleFunc("/sensors/microwave/doorCount", getDoorsStatusTimesEachDay(sendData)).Methods("GET", "OPTIONS")
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

func getMicroWaveSensorData(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.Open("data/10014_FFFFFFFF006eb624_Microwave-Door-Sensor.csv")
		starttime := int64(-1)
		if r.URL.Query().Get("starttime") != "" {
			starttime, _ = strconv.ParseInt(r.URL.Query().Get("starttime"), 10, 64)
		}

		endtime := int64(math.MaxInt64)
		if r.URL.Query().Get("starttime") != "" {
			endtime, _ = strconv.ParseInt(r.URL.Query().Get("endtime"), 10, 64)
		}

		status := ""
		if r.URL.Query().Get("status") != "" {
			status = r.URL.Query().Get("status")
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reader := csv.NewReader(bufio.NewReader(content))
		defer content.Close()

		type Door struct {
			DayofWeek  int    `json:"dayofweek"`
			Day        int    `json:"day"`
			Month      int    `json:"month"`
			Year       int    `json:"year"`
			Doorstatus string `json:"doorstatus"`
			Count      int    `json:"count"`
		}

		var doors []Door
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			}
			timestamp, _ := strconv.ParseInt(line[4], 10, 64)

			if timestamp > int64(starttime) && timestamp < int64(endtime) {
				//fmt.Println(time.Unix(int64(timestamp), 0))

				if status == "" || line[12] == status {
					// fmt.Println(time.Now())
					doors = append(doors, Door{
						DayofWeek:  int(time.Unix(timestamp/1000, timestamp%1000).Weekday()),
						Day:        time.Unix(timestamp/1000, timestamp%1000).Day(),
						Month:      int(time.Unix(timestamp/1000, timestamp%1000).Month()),
						Year:       time.Unix(timestamp/1000, timestamp%1000).Year(),
						Doorstatus: line[12],
						Count:      -1,
					})
				}
			}
		}

		context.Set(r, "objs", doors)
		h.ServeHTTP(w, r)
	})
}

func getDoorsStatusTimesEachDay(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.Open("data/10014_FFFFFFFF006eb624_Microwave-Door-Sensor.csv")
		starttime := int64(-1)
		if r.URL.Query().Get("starttime") != "" {
			starttime, _ = strconv.ParseInt(r.URL.Query().Get("starttime"), 10, 64)
		}

		endtime := int64(math.MaxInt64)
		if r.URL.Query().Get("starttime") != "" {
			endtime, _ = strconv.ParseInt(r.URL.Query().Get("endtime"), 10, 64)
		}

		status := ""
		if r.URL.Query().Get("status") != "" {
			status = r.URL.Query().Get("status")
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reader := csv.NewReader(bufio.NewReader(content))
		defer content.Close()

		type Door struct {
			DayofWeek  string    `json:"dayofweek"`
			Day        string    `json:"day"`
			Month      string    `json:"month"`
			Year       string    `json:"year"`
			Count      int    `json:"count"`
		}

		var doors []Door
		mp := make(map[string]int)

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			}
			timestamp, _ := strconv.ParseInt(line[4], 10, 64)

			if timestamp > int64(starttime) && timestamp < int64(endtime) {
				//fmt.Println(time.Unix(int64(timestamp), 0))
				dayofweek := int(time.Unix(timestamp/1000, timestamp%1000).Weekday())
				day := time.Unix(timestamp/1000, timestamp%1000).Day()
				month := int(time.Unix(timestamp/1000, timestamp%1000).Month())
				year := time.Unix(timestamp/1000, timestamp%1000).Year()
				key := strconv.Itoa(dayofweek) + "/" + strconv.Itoa(day) + "/"  + strconv.Itoa(month) + "/" + strconv.Itoa(year)

				if status == "" || line[12] == status {
					// fmt.Println(time.Now())
					mp[key]++
				}
			}
		}

		for key, element := range mp {
			i := strings.Index(key, "/")
			daysofWeek := key[:i]
			remain := key[i+1:]

			i = strings.Index(remain, "/")
			day := remain[:i]
			remain = remain[i+1:]

			i = strings.Index(remain, "/")
			month := remain[:i]
			year := remain[i+1:]

			doors = append(doors, Door{
				DayofWeek:  daysofWeek,
				Day:        day,
				Month:      month,
				Year:       year,
				Count:      element,
			})
		}

		context.Set(r, "objs", doors)
		h.ServeHTTP(w, r)
	})
}

func sendData(w http.ResponseWriter, r *http.Request) {
	allData := context.GetAll(r)
	objs := allData["objs"]
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(objs)
}
