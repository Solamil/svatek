package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
	"strings"
	"strconv"
	"os"
)

var File string = "nameday_cz_sk.txt"
var PrettyFile string = "nameday_cz_sk_pretty.txt"

var PlainTextAgents = []string{
	"curl",
	"httpie",
	"lwp-request",
	"wget",
	"python-requests",
	"python-httpx",
	"openbsd ftp",
	"powershell",
	"fetch",
	"aiohttp",
	"http_get",
	"xh",
}
//Meniny má
//Dnes má svátek
type todayUrlParams struct {
	Country [1]string `json:"country"`
	Pretty []string	`json:"p"`
}

type easterUrlParams struct {
	Year [1]string `json:"year"`
}

var Country string = "cs-CZ"
const Lines int = 367
const Columns int = 3
var list = [Lines][Columns]string{{}}

func main() {
	file := File
	prepareList(file)
	fmt.Println(velikonoce(2022))
	http.HandleFunc("/index.html", index_handler)
	http.HandleFunc("/today", today_handler)
	http.HandleFunc("/dnes", today_handler)
	http.HandleFunc("/velikonoce", easterday_handler)
	http.HandleFunc("/nameday_cz_sk.txt", file_handler)
	http.HandleFunc("/nameday_cz_sk_pretty.txt", file_handler)
	http.HandleFunc("/", index_handler)
	http.ListenAndServe(":8903", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	agent := strings.Split(r.UserAgent(), "/")[0]
	if index := getIndex(PlainTextAgents, agent); index == -1 {
		http.ServeFile(w, r, "web/index.html")
		return
	}
	t := time.Now()
	answer := fmt.Sprintf("%s|%s", getName(list, t, "cs-CZ"), getName(list, t, "sk-SK"))
	w.Write([]byte(answer))

//	http.ServeFile(w, r, "web/index.html")
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func today_handler(w http.ResponseWriter, r *http.Request) {
	country := Country
	t := time.Now()
//	var answer string = ""

	q, _ := url.PathUnescape(r.URL.RawQuery)
	if len(q) != 0 {
		m, _ := url.ParseQuery(q)
		js, _ := json.Marshal(m)

		var param *todayUrlParams
		json.Unmarshal(js, &param)		
		if len(param.Country[0]) > 0 {
			country = param.Country[0]
		}
		if param.Pretty != nil {
			fmt.Println("hello world")	
		}

	}
	result := getName(list, t, country)
	if result == "" {
		result = getName(list, t, Country)
	}

	w.Write([]byte(result))
}

func easterday_handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	year := t.Year()
	q, _ := url.PathUnescape(r.URL.RawQuery)
	if len(q) != 0 {
		m, _ := url.ParseQuery(q)
		js, _ := json.Marshal(m)

		var param *easterUrlParams
		json.Unmarshal(js, &param)		
		
		if len(param.Year[0]) > 0 {
			y, err := strconv.Atoi(param.Year[0])
			year = y
			if err != nil || year <= 1583 {
				w.Write([]byte("Wrong year. Calculation is for gregorian calendar therefore from 1583 above."))
				return
			}
		}
		
	}
	w.Write([]byte(velikonoce(year)))
}

// credits to https://kalendar.beda.cz/vypocet-velikonocni-nedele-v-ruznych-programovacich-jazycich
func velikonoce(rok int) string {
	if rok <= 1583 {
		return ""
	}
	zlateCislo := (rok % 19) + 1
	julEpakta := (11 * zlateCislo) % 30
	stoleti := int(rok / 100) + 1
	slunecniOprava := int(3 * (stoleti - 16) / 4)
	mesicniOprava := int(8 * (stoleti - 15) / 25)
	epakta := (julEpakta - 10 - slunecniOprava + mesicniOprava) % 30
	if epakta < 0 {
		epakta += 30
	}
	tmp := epakta
	if epakta == 24 || (epakta == 25 && zlateCislo > 11) {
		tmp += 1
	}
	pfm := 0 // Paschal Full Moon
	if tmp < 24 {
		pfm = 44 - tmp
	} else {
		pfm = 74 - tmp
	}

	gregOprava := 10 + slunecniOprava
	denTydnePfm := (rok + (int)(rok / 4) - gregOprava + pfm) % 7
	if denTydnePfm < 0 {
		denTydnePfm += 7
	}
	velNedele := pfm + 7 - denTydnePfm
	if velNedele < 32 {
		return fmt.Sprintf("%d. březen", velNedele)	
	} 
	return fmt.Sprintf("%d. duben", velNedele-31)	
}

func prepareList(name string) {
	for i, line := range readFile(name) {
		d := strings.Split(line, "|")
		for j, value := range d {
			list[i][j] = value
		} 
	}
}
func getName(list [Lines][Columns]string, t time.Time, country string) string {
	var result string = ""
	names := getLineByDate(list, t)
	col := getIndex(list[0][:], country)
	if col > 0 && col < Columns {
		result = names[col]
	}
	return result 
}

func getLineByDate(list [Lines][Columns]string, t time.Time) []string {
	var result []string
	date := fmt.Sprintf("%d.%d", t.Day(), int(t.Month()))
	
	for i := 0; i < Lines; i++ {
		if date == list[i][0] {
			result = list[i][:]
			break
		}
	}
	return result
}

func getDate(list [Lines][Columns]string, name string, col int) string {
	var result string = ""
	for i := 0; i < Lines && col > 0 && col < Columns ; i++ {
		if name == list[i][col] {
			result = list[i][0]
			break
		}
		names := strings.Split(list[i][col],"/")
		
		if j := getIndex(names, name); j != -1 {
			result = list[i][0]
			break
		}
	}
	return result 
}

func readFile(name string) []string {
	file, err := os.Open(name)
	if err != nil {
		fmt.Sprintf("Failed to open %s", name)
		return []string{}
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string 

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	
	return text
}

func getIndex(list []string, value string) int {
	var index int = -1
	for i := 0; i < len(list); i++ {
		if list[i] == value {
			index = i	
			break
		}
	}
	return index
}
