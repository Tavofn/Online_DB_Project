package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
type Songschild struct {
	SongID   string
	Song     string
	SongPath string
}
type Artist struct {
	artist string
}
type Songs struct {
	MySongChild     []Songschild //suggested results
	MySongChildData []Songschild //actual results
}
type dbstruct struct {
	mydb *sql.DB
	user User
}
type playlistIDStru struct {
	PlaylistID string
}
type dbstructIsPlaylistSearch struct {
	mydb dbstruct
}
type playlist struct {
	Plchild []playlistChild
}
type playlistChild struct {
	Playlistname string
	PlaylistID   int
	songID       string
	Song         string
	SongPath     string
}

var tpl *template.Template

var isTrue bool = false
var plIDS playlistIDStru

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
func (mydb *dbstruct) login(w http.ResponseWriter, r *http.Request) {
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

			mydb.user.userid = userID
			fmt.Println(mydb.user.userid)

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

			http.Redirect(w, r, "/", http.StatusSeeOther)

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
func (mydb *dbstruct) searchList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		songs, err2 := mydb.mydb.Query("SELECT title FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		songPath, err2 := mydb.mydb.Query("SELECT mp3_file FROM song ")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		var t Songs

		var songName string
		for songs.Next() && songPath.Next() {
			songs.Scan(&songName)
			temp := Songschild{Song: songName[0:strings.Index(songName, "-")]}
			t.MySongChildData = append(t.MySongChildData, temp)

			temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:]}
			t.MySongChildData = append(t.MySongChildData, temp2)

		}
		// newlist := []string{}
		fmt.Println(mydb.user.name_user)
		if mydb.user.name_user != "" {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

		// tpl.ExecuteTemplate(w, "search.html", nil)
		return
	}
	r.ParseForm()

	fmt.Println("here")
	searchType := r.FormValue("selectedOptionNAME")
	search := r.FormValue("mysearch")
	fmt.Println(string(searchType))
	fmt.Println("before")
	if searchType == "Song Name" {
		fmt.Println("i'm here now")
		songs, err2 := mydb.mydb.Query("SELECT title, mp3_file FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		// songPath, err2 := mydb.mydb.mydb.Query("SELECT  FROM song ")
		// if err2 != nil {
		// 	fmt.Println(songs)
		// }
		var t Songs

		var songName string
		var songPathstr string

		for songs.Next() {
			songs.Scan(&songName, &songPathstr)
			if search == songName[strings.Index(songName, "-")+1:] {
				temp := Songschild{Song: songName, SongPath: songPathstr}
				t.MySongChild = append(t.MySongChild, temp)
				fmt.Println(songPathstr)
			}
		}
		if mydb.user.name_user != "" {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

	} else if searchType == "Genre" {
		songs, err2 := mydb.mydb.Query("SELECT genre FROM artist")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		songPath, err2 := mydb.mydb.Query("SELECT mp3_file FROM song ")
		if err2 != nil {
			fmt.Println(songs)
		}
		count := 0
		var t Songs

		var songName string
		var songPathstr string

		for songs.Next() && songPath.Next() {
			songs.Scan(&songName)
			songPath.Scan(&songPathstr)
			temp := Songschild{Song: songName, SongPath: songPathstr}
			t.MySongChild = append(t.MySongChild, temp)
			count++
		}
		if mydb.user.name_user != "" {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

	} else if searchType == "Artist" {
		songs, err2 := mydb.mydb.Query("SELECT title FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		songPath, err2 := mydb.mydb.Query("SELECT mp3_file FROM song ")
		if err2 != nil {
			fmt.Println(songs)
		}

		var t Songs
		var songName string
		var songPathstr string

		for songs.Next() && songPath.Next() {
			songs.Scan(&songName)
			songPath.Scan(&songPathstr)
			if search == songName[0:strings.Index(songName, "-")] {
				temp := Songschild{Song: songName, SongPath: songPathstr}
				t.MySongChild = append(t.MySongChild, temp)
			}
		}
		if mydb.user.username != "" {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

	}

}

func (mydb *dbstruct) searchListPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		songs, err2 := mydb.mydb.Query("SELECT title FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		songPath, err2 := mydb.mydb.Query("SELECT mp3_file FROM song ")
		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		var t Songs

		var songName string
		for songs.Next() && songPath.Next() {
			songs.Scan(&songName)
			temp := Songschild{Song: songName[0:strings.Index(songName, "-")]}
			t.MySongChildData = append(t.MySongChildData, temp)

			temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:]}
			t.MySongChildData = append(t.MySongChildData, temp2)

		}
		// newlist := []string{}
		if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
			fmt.Println(err)
		}
		// tpl.ExecuteTemplate(w, "search.html", nil)
		return
	}
	r.ParseForm()

	fmt.Println("here")
	searchType := r.FormValue("selectedOptionNAME")
	search := r.FormValue("mysearch")
	fmt.Println(string(searchType))
	fmt.Println("before")
	fmt.Println("songID BEFORE conditional ", r.FormValue("songID"))
	if r.FormValue("songID") == "" {

		if searchType == "Song Name" {
			fmt.Println("i'm here now on song name")
			songsID, err2 := mydb.mydb.Query("SELECT songID FROM song")
			if err2 != nil {
				fmt.Println(songsID)
			}
			defer songsID.Close()
			// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
			// if err2 != nil {
			// 	fmt.Println(songs)
			// }
			var t Songs

			var songID int

			for songsID.Next() {
				songsID.Scan(&songID)
				songDetails, err2 := mydb.mydb.Query("SELECT title, mp3_file FROM song WHERE songID=?", songID)
				var songTitle string
				var songPathstr string
				if err2 != nil {
					fmt.Println(err2)
				}
				for songDetails.Next() {

					songDetails.Scan(&songTitle, &songPathstr)
					fmt.Println("song title from query ", songTitle)
					fmt.Println("song search from input ", search)
					if search == songTitle[strings.Index(songTitle, "-")+1:] {
						temp := Songschild{SongID: strconv.Itoa(songID), Song: songTitle, SongPath: songPathstr}
						t.MySongChild = append(t.MySongChild, temp)
						fmt.Println(songPathstr)
					}

				}

			}
			if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
				fmt.Println(err)
			}
		} else if searchType == "Genre" {

		} else if searchType == "Artist" {

			fmt.Println("i'm here now")
			songsID, err2 := mydb.mydb.Query("SELECT songID FROM song")
			if err2 != nil {
				fmt.Println(songsID)
			}
			defer songsID.Close()
			// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
			// if err2 != nil {
			// 	fmt.Println(songs)
			// }
			var t Songs

			var songID int

			for songsID.Next() {
				songsID.Scan(&songID)
				songDetails, err2 := mydb.mydb.Query("SELECT title, mp3_file FROM song WHERE songID=?", songID)
				var songTitle string
				var songPathstr string
				if err2 != nil {
					fmt.Println(err2)
				}
				for songDetails.Next() {

					songDetails.Scan(&songTitle, &songPathstr)
					if search == songTitle[0:strings.Index(songTitle, "-")] {
						temp := Songschild{SongID: string(songID), Song: songTitle, SongPath: songPathstr}
						t.MySongChild = append(t.MySongChild, temp)
						fmt.Println(songPathstr)
					}

				}

			}
			if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
				fmt.Println(err)
			}
		}

	} else {

		fmt.Println("this is my message from options: ", r.FormValue("songID"))
		songidconv, err := strconv.Atoi(r.FormValue("songID"))
		myquery, err := mydb.mydb.Query("SELECT title, mp3_file FROM SONG WHERE songID=?", songidconv)
		myquery.Next()
		var title string
		var path string
		myquery.Scan(&title, &path)
		fmt.Println("query for songID after click on results ", title, path)
		insert, err := mydb.mydb.Prepare("INSERT INTO `playlist_song` (`playlistID`, `song`, `songpath`) VALUES (?, ?, ?);")
		plvalue, err := strconv.Atoi(plIDS.PlaylistID)
		fmt.Println("my playlist ID ahahhahah ", plvalue)
		res, err := insert.Exec(plvalue, title, path)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func (mydb *dbstruct) uploadsong(w http.ResponseWriter, r *http.Request) {
	isFound := false
	fmt.Println(mydb.user.userid, "userid on open")
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
	fmt.Println(mydb.user.userid, " this is my user id for insert on song")
	res, err := insert.Exec(time.Now(), string(artist_name+"-"+title), 0, 0, "/songs/"+string(artist_name+"-"+title)+".mp3", mydb.user.userid, 1, artistidNUM)
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
func (mydb *dbstruct) home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if mydb.user.username == "" {
			tpl.ExecuteTemplate(w, "home_nologin.html", nil)
			return
		}

		myquery, err2 := mydb.mydb.Query("SELECT Playlist_ID,playlist_name FROM PLAYLIST WHERE UserID=?", mydb.user.userid)

		if err2 != nil {
			fmt.Println(err2)
		}

		var playlistname string
		var playlistid int
		var t playlist

		for myquery.Next() {
			//get playlist name and id
			myquery.Scan(&playlistid, &playlistname)
			myquery, err2 := mydb.mydb.Query("SELECT song,songpath FROM PLAYLIST_SONG WHERE PlaylistID=?", playlistid)
			if err2 != nil {
				fmt.Println("test")
			}
			var songname string
			var songpath string
			songnameList := ""
			songpathList := ""
			for myquery.Next() {
				fmt.Println("here on myquery songlist")
				//add songs under playlist name and id
				myquery.Scan(&songname, &songpath)
				songnameList = songnameList + songname + ","
				songpathList = songpathList + songpath + ","

			}
			temp := playlistChild{Playlistname: playlistname, PlaylistID: playlistid, Song: songnameList, SongPath: songpathList}
			t.Plchild = append(t.Plchild, temp)
		}

		tpl.ExecuteTemplate(w, "home.html", t)
	}
	playlistIDVar := r.FormValue("playlistVals")
	res := strings.Split(playlistIDVar, "]")
	plIDS.PlaylistID = res[0]
	fmt.Println(plIDS.PlaylistID)

	http.Redirect(w, r, "/searchPlaylistSong.html", http.StatusSeeOther)

}
func (mydb *dbstruct) createPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "createplay.html", nil)
		return
	}
	playlistName := r.FormValue("playlistTitle")

	myquery, err := mydb.mydb.Prepare("INSERT INTO `playlist` (`date_created`, `time`, `playlist_name`, `UserID`) VALUES (?, ?, ?, ?);")
	res, err := myquery.Exec(time.Now(), 0, string(playlistName), mydb.user.userid)
	if err != nil {
		fmt.Println(err)
	}
	print(res)
	tpl.ExecuteTemplate(w, "createplay.html", "Playlist Created")
}

func (mydb *dbstruct) addAccountSignUp(w http.ResponseWriter, r *http.Request) {
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

	mydb := dbstruct{mydb: myDB()}
	// test := webHandler{mu: mux}
	mux.HandleFunc("/login.html", mydb.login)
	mux.HandleFunc("/signup", mydb.addAccountSignUp)

	mux.HandleFunc("/", mydb.home)
	mux.HandleFunc("/upload_song", mydb.uploadsong)
	mux.HandleFunc("/search.html", mydb.searchList)

	mux.HandleFunc("/searchPlaylistSong.html", mydb.searchListPlaylist)
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
	mux.HandleFunc("/createplay.html", mydb.createPlaylist)

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images"))))
	mux.Handle("/songs/", http.StripPrefix("/songs/", http.FileServer(http.Dir("./songs"))))
	listenAddr := fmt.Sprintf(":%d", args.port)

	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, mux))

}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine

**test**
*/
