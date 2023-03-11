package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type CommandArgs struct {
	port   int
	webDir string
}

func parseArgs() *CommandArgs {
	var port int
	var webDir string
	flag.IntVar(&port, "port", 8086, "TCP port that the HTTP server will listen on.")
	flag.StringVar(&webDir, "web-dir", "web", "The directory from which files are served over HTTP.")
	flag.Parse()
	return &CommandArgs{
		port:   port,
		webDir: webDir,
	}
}

func main() {
	args := parseArgs()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	http.HandleFunc("/newAccount.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/newAccount.html")
	})

	http.HandleFunc("/forgotpassword.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/forgotpassword.html")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	listenAddr := fmt.Sprintf(":%d", args.port)

	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, nil))

}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine

**test**
*/
