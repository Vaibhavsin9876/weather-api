package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Weather struct {
	ID       int     `json:"-" gorm:"primary_key"`
	Country  string  `json:"country"`
	City     string  `json:"city" sql:"unique"`
	MinTemp  string  `json:"mintemp"`
	MaxTemp  string  `json:"maxtemp"`
	Wind     float32 `json:"wind"`
	Pressure int     `json:"pressure"`
	Humadity int     `json:"humadity"`
}

var Db *gorm.DB
var err error
func Routes() {
	r := mux.NewRouter()
	r.HandleFunc("/", CreateWeather).Methods("POST")
	r.HandleFunc("/{city}", GetWeather).Methods("GET")
	r.HandleFunc("/{id}", DeleteWeather).Methods("DELETE")
	r.HandleFunc("/{id}", UpdateWeather).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", r))
}
func DbSetup() {
	Db, err = gorm.Open("sqlite3", "w.db")
	if err != nil {
		panic(err.Error())
	}
	Db.AutoMigrate(Weather{})
}
func main() {
	DbSetup()
	Routes()
}
func CreateWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var weather Weather
	json.NewDecoder(r.Body).Decode(&weather)
	Db.Create(&weather)
	json.NewEncoder(w).Encode(weather)
}
func GetWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	var weather []Weather
	json.NewDecoder(r.Body).Decode(&weather)
	Db.Where("city=?", params["city"]).Find(&weather)
	json.NewEncoder(w).Encode(weather)
}
func DeleteWeather(w http.ResponseWriter, r *http.Request) {
	var weather Weather
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&weather)
	Db.Delete(&weather, params["id"])
	fmt.Fprintf(w, "weather %s is deleted", params["id"])
}
func UpdateWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	var weather Weather
	Db.Where("id=?", params["id"]).First(&weather)
	json.NewDecoder(r.Body).Decode(&weather)
	Db.Save(&weather)

	json.NewEncoder(w).Encode("update success")
}
