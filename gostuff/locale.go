package gostuff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type IPLocation struct {
	IP           string
	Country_code string
	Country_name string
	Region_code  string
	Region_name  string
	City         string
	Zip_code     string
	Time_zone    string
	Latitude     float32
	Longitude    float32
	Metro_code   string
}

// Sets the country the player is from in the database the first time they register
// returns back what the country name that was set
func setCountry(username string, ipAddress string) string {

	var country = "globe"

	client := TimeOutHttp(5)

	response, err := client.Get("http://api.ipstack.com/" + ipAddress + "?access_key=" +
		ReadOneLine("secret/ipstack.txt") + "&output=json&legacy=1")
	if response == nil {
		fmt.Println("URL time out for http://api.ipstack.com/ in setCountry")
		return "globe"
	}
	defer response.Body.Close()
	if err != nil {
		fmt.Println("error in get language 1", err)
		return "globe"
	}

	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error in get language 2", err)
		return "globe"
	}

	var ipLocation IPLocation

	if err := json.Unmarshal(htmlData, &ipLocation); err != nil {
		fmt.Println("error in get language 3", string(htmlData), err)
		return "globe"
	}

	stmt, err := db.Prepare("UPDATE userinfo SET country=? WHERE username=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println("error in get language 4", err)
		return "globe"
	}
	country = strings.ToLower(ipLocation.Country_code)
	_, err = stmt.Exec(country, username)
	if err != nil {
		fmt.Println("error in get language 5", err)
		return "globe"
	}
	return country
}

// Fetches country from database for a given player every time they login
// If country is null then it returns blank string which should be checked
func GetCountry(username string) string {

	//globe.png is default country flag
	var country = "globe"

	err := db.QueryRow("SELECT country from userinfo WHERE username=?", username).Scan(&country)
	if err != nil { // then country is nil
		//fmt.Println("Error getting country for", username, err)
		return "globe"
	}
	return country
}
