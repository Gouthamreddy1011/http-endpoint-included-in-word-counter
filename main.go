package main

import (
	"encoding/json"
	"fmt"       //formatting package
	"io/ioutil" //i/outility functions
	"log"       //implements simple logging package
	"net/http"  //http provides client and server implementation get,head,post,and post form
	"regexp"    //regular expression search
	"sort"      //premitives for sorting slices and user defined collections
	"strconv"   //conversions to and from string representations of data type
	"strings"   //manuplates uft-8 encoded strings
)

var userdata []string

type WordCount struct { // to count number of repeating words
	word  string // words in the form string
	count int    //int for integers
}

func inputHandler(w http.ResponseWriter, r *http.Request) { //not found replies to request with an http 404 not found error

	body, err := ioutil.ReadAll(r.Body) //reading

	if err != nil { //error nil handling
		log.Fatal(err) //if error
	}
	getdata := string(body) // data body of string

	re, err := regexp.Compile(`[^\w]`) //re compile if matches any single letter,number, underscore
	if err != nil {                    //error nil handling
		log.Fatal(err) // if error
	}

	getdata = re.ReplaceAllString(getdata, " ") // replacing strings
	getdata = strings.ToLower(getdata)          //pass strings as argument to getdata

	str, err := json.Marshal(getdata) //  marshal/unmarshal
	if err != nil {                   //  if error
		fmt.Printf("Error: %s", err.Error()) // error string
	}

	userinput := string(str[:]) //user input strings

	words := strings.Fields(userinput) // user input words

	m := make(map[string]int) //including word counter
	for _, word := range words {

		_, ok := m[word]
		if !ok {
			m[word] = 1
		} else {
			m[word]++
		}
	}
	wordCounts := make([]WordCount, 0, len(m))
	for key, val := range m {
		wordCounts = append(wordCounts, WordCount{word: key, count: val}) //appending slices
	}
	sort.Slice(wordCounts, func(i, j int) bool { // sorting slice function i,j
		return wordCounts[i].count > wordCounts[j].count // i>j
	})
	for i := 0; i < len(wordCounts) && i < 10; i++ { // for statement post incrementation
		userdata = append(userdata, " word: ", wordCounts[i].word, " occurs: ", strconv.Itoa(wordCounts[i].count))
	}
	fmt.Fprintf(w, "The count is: %s", userdata) //stores in userdata
}

func main() {

	http.HandleFunc("/input", inputHandler)

	err := http.ListenAndServe(":8000", nil) //server
	if err != nil {                          //if no error
		log.Fatal("ListenAndServe: ", err) // if error
	}

}
