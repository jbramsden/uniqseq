//Uniseq example for a working URL shortner.
// This will store the shortern URL in memory, so if the program restarts all URLs are lost!
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/jbramsden/uniqseq"
)

//Global value for uniqseq
var a *uniqseq.UniqueString

//Mdb - Storage for Shortern URL (Would be better to store in DB or Redis or something similar
var Mdb = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

//frontpage - Displays a form to capture the URL to shortern or if short URL passed it will redirect the brower to the full url.
//            This just covers the happy path, error handling should be added for production use(plus handling for favicon.ico)
func frontpage(w http.ResponseWriter, r *http.Request) {
	//Display the html page with form to shortern URL
	if len(r.URL.Path) == 1 {
		io.WriteString(w, "<html><head><link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css' integrity='sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ' crossorigin='anonymous'></head><body><div class='row h-100'><div class='col-sm-12 my-auto'><div class='w-50 mx-auto'><form action='/post' method='post' class='form-horizontal'><fieldset><div class='form-group'> <div class='col-md-12'><input id='url' name='url' type='text' placeholder='http://github.com/jbramsden' class='form-control input-md'><span class='help-block'>Enter a URL </span></div><div class='col-md-12'><input type='submit' class='btn btn-success'></div></div></fieldset></form></div></div></div></body></html>")
	} else {
		//find the URL and create the HTML for redirection
		Mdb.RLock()
		rurl := Mdb.m[r.URL.Path]
		Mdb.RUnlock()
		redirectHTML := fmt.Sprintf("<html><head><meta http-equiv='refresh' content='0; url=%s'/></head></html>", rurl)
		io.WriteString(w, redirectHTML)
	}
}

//handlePost - This take the post for the URL to shorten and gets the next uniseq number and adds it to the map along with the URL.
//             To improve the posted url, this needs to be checked to validate that an actual url is passed in the form.
func handlePost(w http.ResponseWriter, r *http.Request) {
	var H string
	if err := r.ParseForm(); err != nil {
		H = fmt.Sprintf("<html><head><link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css' integrity='sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ' crossorigin='anonymous'></head><body><div class='row h-100'><div class='col-sm-12 my-auto'><div class='w-50 mx-auto'><h2>Error: URL must be provided</h2></div></div></div></body></html>")
	} else {

		url := r.PostFormValue("url")
		//Get the next sequence number from UniSeq
		u, _ := a.Next()

		//Add a forward slash in front of it before storing. Makes it easier to find as r.URL.Path does not strip it
		us := fmt.Sprintf("/%s", u)

		//Create a lock and then update the map with the uniseq and url
		Mdb.Lock()
		Mdb.m[us] = url
		Mdb.Unlock()

		H = fmt.Sprintf("<html><head><link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css' integrity='sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ' crossorigin='anonymous'></head><body><div class='row h-100'><div class='col-sm-12 my-auto'><div class='w-50 mx-auto'>Short URL :<h2>%s%s</h2></div></div></div></body></html>", r.Host, us)
	}
	io.WriteString(w, H)
}

func main() {
	//Get the pointer with default value set
	a = uniqseq.Create()
	//The following will generate the following unique sequence AIIII xIIII 4IIII or somthing similar as Jumbler is active
	a.CharacterSet = uniqseq.NoVowels
	a.Jumbler = true
	a.StartLength = 5
	a.BlankFillChar = "I"
	a.LastCharInc = false
	//Initialise so the sequence is at the begining and also that the character set are jumbled.
	a.Init()

	//Frontpage displays the form to create a URL and it also handles the redirect for short URLs
	http.HandleFunc("/", frontpage)
	//Post handles the post action from the form and generate a short URL
	http.HandleFunc("/post", handlePost)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
