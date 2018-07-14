package main

import (
	"fmt"
	"github.com/brucemaclin/urlshorter"
	"io"
	"net/http"
	"os"
)

func getShorterHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	origURL := req.FormValue("url")
	if origURL == "" {
		io.WriteString(w, "need url to shorter")
		return
	}
	//fmt.Println("origURL:", origURL)
	shortURL, err := shorter.ShorterURL(origURL)
	if err != nil {
		io.WriteString(w, "fail to get shortURL:"+err.Error())
		return
	}
	io.WriteString(w, "short:t.co/"+shortURL)
	return
}
func main() {

	db := &shorter.DefaultDB{}
	db.Init()
	shorter.InitWithDB(db)
	http.HandleFunc("/shorter", getShorterHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("fail to listen and server:", err)
		os.Exit(-1)
	}

}
