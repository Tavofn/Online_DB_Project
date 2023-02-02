package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type CommandArgs struct {
	port   int
	webDir string
}

func parseArgs() *CommandArgs {
	var port int
	var webDir string
	flag.IntVar(&port, "port", 8085, "TCP port that the HTTP server will listen on.")
	flag.StringVar(&webDir, "web-dir", "web", "The directory from which files are served over HTTP.")
	flag.Parse()
	return &CommandArgs{
		port:   port,
		webDir: webDir,
	}
}

type handler struct {
	content string
}
type myvals struct {
	Http10 string
	Http11 string
	Http12 string
	Http13 string

	Tls10 string
	Tls11 string
	Tls12 string
	Tls13 string
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, err := template.ParseFiles("./web/index.html")
	if err != nil {
		fmt.Print("iuykfsduyfsdf")
	}
	t.Execute(w, nil)
}

func main() {
	args := parseArgs()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	listenAddr := fmt.Sprintf(":%d", args.port)
	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, nil))
}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine
*/
