package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tcolgate/mp3"
)

type CommandArgs struct {
	port   int
	webDir string
}

type webHandler struct {
	mu *http.ServeMux
}

type User struct {
	userid      int
	username    string
	password    string
	email       string
	date_regist string
	name_user   string
}

var tpl *template.Template

var isTrue bool = false

func myDB() *sql.DB {
	//pass, err := ioutil.ReadFile("pass.txt")
	//user, err := ioutil.ReadFile("username.txt")
	db, err := sql.Open("mysql", "hpalma:5802**@tcp(team3-music-database-2023.mysql.database.azure.com:3306)/3380-project?tls=skip-verify")
	if err != nil {
		panic(err)
	}
	return db
}
func (we webHandler) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("im here")
	if r.Method == "GET" {
		fmt.Println("in get rn")
		err := tpl.ExecuteTemplate(w, "login.html", "")
		if err != nil {
			fmt.Println("on GET MY ERROR")
			fmt.Println(err)
		}
		return
	}

	// fmt.Println("out of get")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	query, err := myDB().Query("SELECT username,password FROM USER")
	// fmt.Println("out of get on query")
	if err != nil {
		fmt.Println("query fail")
		fmt.Println(err.Error())
	}
	defer query.Close()
	fmt.Println("out of get after defer close")

	if len(username) == 0 || len(password) == 0 {
		fmt.Println("empty")
		err2 := tpl.ExecuteTemplate(w, "login.html", "please Fill All Spaces")
		if err2 != nil {
			fmt.Println("on EMPTY ERROR")
			fmt.Println(err)
		}
		return
	}
	for query.Next() {
		var t User
		query.Scan(&t.username, &t.password)
		if t.username == username && t.password == password {
			fmt.Print("correct")

			err := tpl.ExecuteTemplate(w, "home.html", "")
			if err != nil {
				fmt.Println("ON HOME ERROR")
				fmt.Println(err)
			}
			return
		}

	}
	fmt.Println("incorrect")
	err2 := tpl.ExecuteTemplate(w, "login.html", "Login failed, please try again")
	if err2 != nil {
		fmt.Println("on login failed")
		fmt.Println(err)
	}
	return
	//test
}
func uploadsong(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "upload_song.html", nil)
		return
	}
	r.ParseForm()
	fmt.Println("sfdsdfsdf first time on post")
	title := r.FormValue("song_name")

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("songs", title+".mp3")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	t := 0.0

	if err != nil {
		fmt.Println(err)
		return
	}

	d := mp3.NewDecoder(tempFile)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}

		t = t + f.Duration().Seconds()
	}

	fmt.Println(t)

	insert, err2 := myDB().Prepare("INSERT INTO `song` (`release_date`, `title`, `time`, `average_rating`, `mp3_file`, `listens`) VALUES (?, ?, ?, ?, ?, '1');")
	if err2 != nil {
		fmt.Println(err2)
	}
	res, err := insert.Exec(time.Now().UTC(), string(title), t, 0, "/songs/"+title, 0)
	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "upload_song.html", "Error inserting data, please check all fields.")
		return
	}
	tpl.ExecuteTemplate(w, "upload_song.html", nil)
}

func addAccountSignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****insertHandler running*****")
	if r.Method == "GET" {
		fmt.Println("here")
		tpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}
	r.ParseForm()
	fmt.Println("sfdsdfsdf first time on post")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	email := r.FormValue("email")
	username := r.FormValue("username")
	passw := r.FormValue("password")
	passwC := r.FormValue("CPassword")

	if fname == "" || lname == "" || email == "" || passw == "" || username == "" || passwC == "" {
		tpl.ExecuteTemplate(w, "signup.html", "Please fill out entire form")
		return
	}

	query, err := myDB().Query("SELECT email FROM USER")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer query.Close()

	for query.Next() {
		var t User
		query.Scan(&t.email)
		if t.email == email {
			tpl.ExecuteTemplate(w, "signup.html", "This Account already exists")
			return
		}
	}

	queryUName, err := myDB().Query("SELECT username FROM USER")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer queryUName.Close()

	for queryUName.Next() {
		var t User
		queryUName.Scan(&t.username)
		if t.username == username {
			tpl.ExecuteTemplate(w, "signup.html", "This username already exists")
			return
		}
	}

	insert, err2 := myDB().Prepare("INSERT INTO `user` (`username`, `password`, `email`, `date_registered`, `name_of_user`, `access_level`) VALUES (?, ?, ?, ?, ?, '1');")
	//(`username`, `password`, `date_registered`, `name_of_user`, `access_level`)
	if err2 != nil {
		fmt.Println("error on prep")
		panic(err2.Error())
	} else {
		fmt.Println("prepare completed")
	}

	res, err := insert.Exec(username, passw, email, time.Now().UTC(), string(fname+"_"+lname))

	if err != nil {
		fmt.Println("error on insert")
		panic(err.Error())
	}

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "signup.html", "Error inserting data, please check all fields.")
		return
	}
	tpl.ExecuteTemplate(w, "signup.html", "Product Successfully Inserted")
	// defer insert.Close()

	fmt.Println("sdkjfhsdf here")

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
	tpl, _ = template.ParseGlob("./web/*.html")

	fmt.Println("Successful Connection to Database!")
	args := parseArgs()

	// defer db.Close()
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web/login.html")
	// })
	// http.HandleFunc("/signup.html", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web/signup.html")
	// })
	mux := http.NewServeMux()
	test := webHandler{mu: mux}
	mux.HandleFunc("/", test.login)
	mux.HandleFunc("/signup", addAccountSignUp)

	mux.HandleFunc("/home.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/home.html")
	})
	mux.HandleFunc("/upload_song", uploadsong)
	mux.HandleFunc("/search.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/search.html")
	})
	mux.HandleFunc("/forgotpassword.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/forgotpassword.html")
	})
	mux.HandleFunc("/myaccount.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/myaccount.html")
	})
	mux.HandleFunc("/AO.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/myaccount.html")
	})
	mux.HandleFunc("/createplay.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/createplay.html")
	})

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images"))))
	listenAddr := fmt.Sprintf(":%d", args.port)

	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, mux))

}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine

**test**
*/
