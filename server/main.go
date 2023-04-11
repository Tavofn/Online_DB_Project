package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
type Songs struct {
	Song         string
	SonglistData []string
	Songlist     []string
	Count        int
}
type Artist struct {
	artist string
}
type SongsChild struct {
	Song        string
	Songlist    []string
	Count       int
	MySongChild Songs
}
type dbstruct struct {
	mydb *sql.DB
	user User
}

var tpl *template.Template

var isTrue bool = false

/*
TODO

*full functioning music player
**with UI interface
*playlist functionality
**front page, search through foreign key (userID) linked to songs
*foreign key for artist inside song
*




*/

func myDB() *sql.DB {
	//pass, err := ioutil.ReadFile("pass.txt")
	//user, err := ioutil.ReadFile("username.txt")
	db, err := sql.Open("mysql", "hpalma:5802**@tcp(team3-music-database-2023.mysql.database.azure.com:3306)/3380-project?tls=skip-verify")
	if err != nil {
		panic(err)
	}
	return db
}
func (mydb dbstruct) login(w http.ResponseWriter, r *http.Request) {
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

	query, err := mydb.mydb.Query("SELECT username,password FROM USER")
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
			fmt.Println("correct")

			mydb.user.username = username
			mydb.user.password = password
			myquery := fmt.Sprintf("SELECT UserID FROM USER WHERE username=\"%s\"", string(username))
			queryuserID, err := mydb.mydb.Query(myquery)
			queryEmail, err := mydb.mydb.Query("SELECT email FROM USER WHERE username=?", string(username))
			myqueryDate := fmt.Sprintf("SELECT date_registered FROM USER WHERE username=\"%s\"", username)
			queryDate, err := mydb.mydb.Query(myqueryDate)
			if err != nil {
				fmt.Println(err, "error here")
			} else {
				fmt.Println("no error on queries login")
			}
			userID := -1
			queryuserID.Next()
			queryuserID.Scan(&userID)
			fmt.Println(userID)
			mydb.user.userid = userID

			var email string
			queryEmail.Next()
			if err := queryEmail.Scan(&email); err != nil {
				panic(err)
			}
			fmt.Println(email)
			mydb.user.email = email

			date := time.Now().Format("2006-01-02 15:04:05")
			queryDate.Next()
			queryDate.Scan(&date)
			fmt.Println(date)
			mydb.user.date_regist = date

			err2 := tpl.ExecuteTemplate(w, "home.html", "")
			if err2 != nil {
				fmt.Println("ON HOME ERROR")
				fmt.Println(err2)
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
func (mydb dbstruct) searchList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		songs, err2 := mydb.mydb.Query("SELECT title FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		count := 0
		var t Songs
		list := []string{}
		for songs.Next() {
			songs.Scan(&t.Song)
			list = append(list, t.Song)
			count++
		}
		newlist := []string{}

		for x := 0; x < len(list)-1; x++ {
			newlist = append(newlist, list[x][0:strings.Index(list[x], "-")])
			// fmt.Println(strings.Index(list[x], "-"))
			// fmt.Println(list[x][0:strings.Index(list[x], "-")])
		}
		for x := 0; x < len(list); x++ {
			newlist = append(newlist, list[x][strings.Index(list[x], "-")+1:len(list[x])])
			// fmt.Println(list[x][strings.Index(list[x], "-")+1 : len(list[x])])
		}

		t.SonglistData = newlist
		t.Count = 10
		if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
			fmt.Println(err)
		}
		// tpl.ExecuteTemplate(w, "search.html", nil)
		return
	}
	r.ParseForm()

	// fmt.Println("here")
	// searchType := r.FormValue("selectedOptionNAME")
	// search := r.FormValue("mysearch")
	// fmt.Println(string(searchType))
	// fmt.Println("before")
	// if searchType == "Song Name" {
	// 	fmt.Println("i'm here now")
	// 	songs, err2 := mydb.mydb.Query("SELECT title FROM song")
	// 	if err2 != nil {
	// 		fmt.Println(songs)
	// 	}
	// 	defer songs.Close()
	// 	count := 0
	// 	var t Songs
	// 	list := []string{}
	// 	for songs.Next() {
	// 		songs.Scan(&t.Song)
	// 		if strings.Index(t.Song[0:strings.Index(t.Song, "-")-1], search) != -1 {
	// 			fmt.Println("found")
	// 			list = append(list, t.Song)
	// 		}
	// 		fmt.Println(t.Song)
	// 		count++
	// 	}
	// 	t.Songlist = list
	// 	t.Count = 10
	// 	if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
	// 		fmt.Println(err)
	// 	}
	// } else if searchType == "Genre" {
	// 	songs, err2 := mydb.mydb.Query("SELECT genre FROM artist")
	// 	if err2 != nil {
	// 		fmt.Println(songs)
	// 	}
	// 	defer songs.Close()
	// 	count := 0
	// 	var t Songs
	// 	list := []string{}
	// 	for songs.Next() {
	// 		songs.Scan(&t.Song)

	// 		list = append(list, t.Song)
	// 		count++
	// 	}
	// 	t.Songlist = list
	// 	t.Count = 10
	// 	if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
	// 		fmt.Println(err)
	// 	}
	// } else if searchType == "Artist" {
	// 	songs, err2 := mydb.mydb.Query("SELECT title FROM song")
	// 	if err2 != nil {
	// 		fmt.Println(songs)
	// 	}
	// 	defer songs.Close()
	// 	count := 0
	// 	var t Songs
	// 	list := []string{}
	// 	for songs.Next() {
	// 		songs.Scan(&t.Song)
	// 		if t.Song == search {
	// 			list = append(list, t.Song[0:strings.Index(t.Song, "-")])
	// 		}
	// 		count++
	// 	}
	// 	t.Songlist = list
	// 	t.Count = 10
	// 	if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

}

func (mydb dbstruct) uploadsong(w http.ResponseWriter, r *http.Request) {
	isFound := false
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "upload_song.html", nil)
		return
	}

	r.ParseForm()
	fmt.Println("sfdsdfsdf first time on post")

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	tempFile, err := os.Create("")
	title := r.FormValue("song_name")
	artist_name := r.FormValue("artist_name")
	if title == "" {
		tempFile, err = os.Create("songs/" + handler.Filename)
	} else {
		tempFile, err = os.Create("songs/" + artist_name + "-" + title + ".mp3")

	}
	fmt.Println(title + " this is my tite")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// t := 0.0

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// d := mp3.NewDecoder(tempFile)
	// var f mp3.Frame
	// skipped := 0

	// for {

	// 	if err := d.Decode(&f, &skipped); err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	t = t + f.Duration().Seconds()
	// }

	// fmt.Println(t)
	// userid, err3 := mydb.mydb.Prepare("SELECT UserID FROM user WHERE username = AND password = ")
	searchArtist, err2 := mydb.mydb.Query("SELECT artist_name FROM artist")
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("no error on search artist")
	}

	for searchArtist.Next() {
		var t Artist
		searchArtist.Scan(&t.artist)
		if t.artist == r.FormValue("artist_name") {
			isFound = true
		}
	}
	fmt.Println("after search artist next")
	if !isFound {
		insertARTIST, err2 := mydb.mydb.Prepare("INSERT INTO `artist` (`artist_name`, `genre`, `average_rating`) VALUES (?, ?, ?);")
		fmt.Println("after prep")
		res2, err2 := insertARTIST.Exec(string(artist_name), string(r.FormValue("genre")), 0)
		fmt.Println("after execute")

		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println("no error")
		}
		fmt.Println(res2)
	}
	artistname := string(r.FormValue("artist_name"))
	artistIDQuery := fmt.Sprintf("SELECT ArtistID FROM ARTIST WHERE artist_name=\"%s\"", artistname)
	myartistID, err2 := mydb.mydb.Query(artistIDQuery)
	fmt.Println("after query artistid")

	myartistID.Next()
	artistidNUM := 0
	myartistID.Scan(&artistidNUM)

	fmt.Println("after scan")
	insert, err2 := mydb.mydb.Prepare("INSERT INTO `song` (`release_date`, `title`, `time`, `average_rating`, `mp3_file`, `UserID`, `listens`, `artist_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	fmt.Println("after prep insert song")

	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("no error")
	}

	res, err := insert.Exec(time.Now(), string(artist_name+"-"+title), 0, 0, "/songs/"+string(artist_name+"-"+title), mydb.user.userid, 1, artistidNUM)
	fmt.Println("after song exe")

	// rowsAffec, _ := res.RowsAffected()
	fmt.Println(res)

	if err != nil {
		fmt.Println("Error inserting row:", err)
		// tpl.ExecuteTemplate(w, "upload_song.html", "Error inserting data, please check all fields.")
		return
	}

	tpl.ExecuteTemplate(w, "upload_song.html", "Song successfuly uploaded")
}
func (mydb dbstruct) home(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "home_nologin.html", nil)
}
func (mydb dbstruct) addAccountSignUp(w http.ResponseWriter, r *http.Request) {
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

	mydb := dbstruct{mydb: myDB(), user: User{}}
	// test := webHandler{mu: mux}
	mux.HandleFunc("/login.html", mydb.login)
	mux.HandleFunc("/signup", mydb.addAccountSignUp)

	mux.HandleFunc("/", mydb.home)
	mux.HandleFunc("/upload_song", mydb.uploadsong)
	mux.HandleFunc("/search.html", mydb.searchList)
	mux.HandleFunc("/forgotpassword.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/forgotpassword.html")
	})
	mux.HandleFunc("/myaccount.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/myaccount.html")
	})
	mux.HandleFunc("/editP.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/editP.html")
	})
	mux.HandleFunc("/reports.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/reports.html")
	})
	mux.HandleFunc("/changepass.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/changepass.html")
	})
	mux.HandleFunc("/AO.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/AO.html")
	})
	mux.HandleFunc("/createplay.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/createplay.html")
	})

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images"))))
	mux.Handle("/savedSongs/", http.StripPrefix("/savedSongs/", http.FileServer(http.Dir("./songs"))))
	listenAddr := fmt.Sprintf(":%d", args.port)

	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, mux))

}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine

**test**
*/
