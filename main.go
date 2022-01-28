package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func jsonResponse(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	p, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(p))
}

func messageResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	type Message struct {
		Message string `json:"message"`
	}
	msg := Message{message}
	p, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(p))
}

func handleCounter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		messageResponse(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	countPerWord := map[string]int{}
	words := strings.Fields(string(body))
	for _, word := range words {
		if val, ok := countPerWord[word]; ok {
			countPerWord[word] = val + 1
		} else {
			countPerWord[word] = 1
		}
	}
	topWords := map[string]int{}
	for {
		if len(countPerWord) == 0 || len(topWords) == 10 {
			break
		}
		maxK := ""
		maxV := 0
		for k, v := range countPerWord {
			if v > maxV {
				maxK = k
				maxV = v
			}
		}
		topWords[maxK] = maxV
		delete(countPerWord, maxK)
	}
	jsonResponse(w, topWords, http.StatusOK)
}

func main() {
	router := httprouter.New()
	router.POST("/counter", handleCounter)
	log.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
