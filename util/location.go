package util

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
  "../models"
)

const apiKey = "AIzaSyAkrnaYlsz_-EkOsb3qGM4GP13NXe6BFO4";

var location struct {
	Results []map[string]map[string]interface{}
}

//GetCoordsForAddress .
func GetCoordsForAddress(address string) (models.Location, error) {
	fmt.Println(address)
	res, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + url.QueryEscape(address) + "&key=" + apiKey)
	byteArray, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(byteArray, &location)
	locationMap := location.Results[0]["geometry"]["location"].(map[string]interface{})
	var location models.Location
	location.Lat = locationMap["lat"].(float64)
	location.Lng = locationMap["lng"].(float64)
	return location, err
}