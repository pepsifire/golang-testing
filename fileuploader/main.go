package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

//Config stuff
type Config struct {
	URL    string `json:"URL"`
	KEY    string `json:"KEY"`
	LENGTH int    `json:"FileNameLength"`
	PORT   string `json:"PORT"`
}

var c Config

//LoadConfig : Loads the json file
func LoadConfig() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error: Failed to load config!", "\n", err)
		return
	}
	fmt.Println("Loaded configuration!")
}

//Upload stuff
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		fmt.Fprintf(w, "%v", "Error")
	} else {
		if r.FormValue("key") != c.KEY {
			keyerror := errors.New("Incorrect key")
			fmt.Fprintf(w, "%v", keyerror)
			fmt.Println(keyerror)
			return
		}
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		extension := strings.Split(handler.Filename, ".") //Split filename
		filename := Randomizer(c.LENGTH)
		f, err := os.OpenFile("./test/"+filename+"."+extension[1], os.O_WRONLY|os.O_CREATE, 0666)
		fmt.Fprintf(w, "%v", c.URL+filepath.Dir(f.Name())+"/"+filename+"."+extension[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

//Randomizer : Random string
func Randomizer(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	LoadConfig()
	router := mux.NewRouter()
	router.HandleFunc("/upload", Upload).Methods("POST")
	router.HandleFunc("/", Upload).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+c.PORT, router))
}
