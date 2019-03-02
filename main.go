package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type typeCities struct {
	Requestid int `json:"requestid,omitempty"`
	Cities    []struct {
		Tm   string `json:"tm"`
		Name string `json:"name"`
	} `json:"cities"`
}

func (tc *typeCities) AddID() error {

	id, errRead := ReadDBGetID()

	if errRead != nil {
		return fmt.Errorf("AddID__errRead %s", errRead)
	}

	tc.Requestid = id

	return nil
}

func convertTime(t string) time.Time {
	var tm time.Time

	if strings.Contains(t, ":") {
		tm, _ = time.Parse("2006-01-02 15:04:05", t)
	} else {
		tm, _ = time.Parse("2006-01-02", t)
	}

	return tm
}

func handleSort(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
	w.Header().Set("Content-Type", "application/json")

	body, errReq := ioutil.ReadAll(r.Body)
	if errReq != nil {
		fmt.Fprintf(w, errReq.Error())
	}
	defer r.Body.Close()

	var getReqCities typeCities

	errJSON := json.Unmarshal(body, &getReqCities)
	if errJSON != nil {
		fmt.Fprintf(w, errJSON.Error())
	}

	for j := 0; j < len(getReqCities.Cities)-1; j++ {
		for i := 0; i < len(getReqCities.Cities)-1; i++ {
			if convertTime(getReqCities.Cities[i].Tm).After(convertTime(getReqCities.Cities[i+1].Tm)) {
				f := getReqCities.Cities[i]
				getReqCities.Cities[i] = getReqCities.Cities[i+1]
				getReqCities.Cities[i+1] = f
			}
		}
	}

	for j := 0; j < len(getReqCities.Cities)-1; j++ {
		for i := 0; i < len(getReqCities.Cities)-1; i++ {
			if getReqCities.Cities[i].Name > getReqCities.Cities[i+1].Name {
				f := getReqCities.Cities[i]
				getReqCities.Cities[i] = getReqCities.Cities[i+1]
				getReqCities.Cities[i+1] = f
			}
		}
	}

	jsonSort, _ := json.Marshal(getReqCities)

	errWrite := WriteDB(jsonSort)

	if errWrite != nil {
		fmt.Fprintf(w, errWrite.Error())
	}

	errAdd := getReqCities.AddID()
	if errAdd != nil {
		fmt.Fprintf(w, errAdd.Error())
	}

	resp, _ := json.Marshal(getReqCities)

	w.Write(resp)
}

func handleGetResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	reqId := r.FormValue("requestid")

	body, errReq := ReadDBGetResult(reqId)
	if errReq != nil {
		fmt.Fprintf(w, errReq.Error())
	}

	var getDBCities typeCities

	errJSON := json.Unmarshal(body, &getDBCities)
	if errJSON != nil {
		fmt.Fprintf(w, errJSON.Error())
	}

	resp, _ := json.Marshal(getDBCities)

	w.Write(resp)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/sort", handleSort).Methods("POST")
	r.HandleFunc("/api/getresult", handleGetResult).Methods("GET")

	log.Fatal(http.ListenAndServe(":82", r))

}
