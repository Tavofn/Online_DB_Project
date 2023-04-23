package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
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
type UserReport struct {
	Userid        int
	Username      string
	Password      string
	Email         string
	Date_regist   string
	Name_user     string
	Playlist_name string
	Playlist_id   string
	Uploadedby    string
	Songcounter   string
}
type Songschild struct {
	SongID       int
	Release_date string
	Song         string
	SongPath     string
	UserID       int
	Listens      int
	ArtistID     int
	Genre        string
	Uploadedby   string
}
type Artist struct {
	artist string
}
type Songs struct {
	MySongChild      []Songschild //suggested results
	MySongChildData  []Songschild //actual results
	SelectedSong     string
	SelectedSongPath string
}

type dbstruct struct {
	mydb       *sql.DB
	user       User
	UserReport []UserReport
	Mysongs    Songs
	playlistID string
}
type playlistIDStru struct {
	PlaylistID string
}
type dbstructIsPlaylistSearch struct {
	mydb dbstruct
}
type playlist struct {
	Plchild  []playlistChild
	Top5     []Songschild
	Username string
}
type playlistChild struct {
	Playlistname string
	PlaylistID   int
	songID       string
	Song         string
	SongPath     string
	Listens      int
	Username     string
}
type userListCHILD struct {
	Username string
	UserID   int
}
type userList struct {
	Users []userListCHILD
}
type HolyGrail struct {
	MySongChild      []Songschild
	MyUserReports    []UserReport
	MyPlaylistChilds []playlistChild
	Message          string
}

var tpl *template.Template

var isTrue bool = false
var plIDS playlistIDStru

var store = sessions.NewCookieStore([]byte("super-secret-password"))

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
	pass, err := ioutil.ReadFile("pass.txt")

	//user, err := ioutil.ReadFile("username.txt")
	db, err := sql.Open("mysql", "hpalma:"+string(pass)+"@tcp(team3-music-database-2023.mysql.database.azure.com:3306)/3380-project?tls=skip-verify")
	if err != nil {
		panic(err)
	}

	// defer db.Close()
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
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

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

			myquery := fmt.Sprintf("SELECT UserID FROM USER WHERE username=\"%s\"", string(username))
			queryuserID, err := mydb.mydb.Query(myquery)
			queryEmail, err := mydb.mydb.Query("SELECT email FROM USER WHERE username=?", string(username))
			myqueryDate := fmt.Sprintf("SELECT date_registered FROM USER WHERE username=\"%s\"", username)
			queryDate, err := mydb.mydb.Query(myqueryDate)
			defer queryDate.Close()
			defer queryEmail.Close()
			defer queryuserID.Close()
			if err != nil {
				fmt.Println(err, "error here")
			} else {
				fmt.Println("no error on queries login")
			}
			userID := -1
			queryuserID.Next()
			queryuserID.Scan(&userID)

			// myuserID = userID
			session, _ := store.Get(r, "session")
			session.Values["userID"] = userID

			var email string
			queryEmail.Next()
			if err := queryEmail.Scan(&email); err != nil {
				panic(err)
			}
			session.Values["email"] = email
			fmt.Println(email)
			// mydb.user.email = email
			session.Values["username"] = username
			date := time.Now().Format("2006-01-02 15:04:05")
			queryDate.Next()
			queryDate.Scan(&date)
			fmt.Println(date)
			session.Values["date"] = date
			// mydb.user.date_regist = date
			session.Save(r, w)
			if t.username == "Admin" {
				http.Redirect(w, r, "/Admin.html", http.StatusSeeOther)

			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)

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
func (mydb *dbstruct) searchList(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	if r.Method == "GET" {
		songs, err2 := mydb.mydb.Query("SELECT title FROM song")
		if err2 != nil {
			fmt.Println(songs)
		}
		artist, err2 := mydb.mydb.Query("SELECT artist_name FROM artist")
		if err2 != nil {
			fmt.Println(artist)
		}
		defer songs.Close()
		defer artist.Close()
		var t Songs

		var songName string

		for songs.Next() {
			songs.Scan(&songName)
			temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:] + ":Song"}
			t.MySongChildData = append(t.MySongChildData, temp2)

		}
		var artistname string
		for artist.Next() {
			artist.Scan(&artistname)
			temp := Songschild{Song: artistname + ":Artist"}
			t.MySongChildData = append(t.MySongChildData, temp)

		}
		x := 0
		genrelist := [6]string{"Hip-Hop", "Pop", "Rock", "Country", "Classical", "Jazz"}
		for x < 6 {
			temp := Songschild{Song: genrelist[x] + ":Genre"}
			t.MySongChildData = append(t.MySongChildData, temp)
			x++
		}

		// newlist := []string{}

		if ok {
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

	songs, err2 := mydb.mydb.Query("SELECT title FROM song")
	if err2 != nil {
		fmt.Println(songs)
	}
	artist, err2 := mydb.mydb.Query("SELECT artist_name FROM artist")
	if err2 != nil {
		fmt.Println(artist)
	}
	defer songs.Close()
	defer artist.Close()
	var t Songs

	var songName string

	for songs.Next() {
		songs.Scan(&songName)
		temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:] + ":Song"}
		t.MySongChildData = append(t.MySongChildData, temp2)

	}
	var artistname string
	for artist.Next() {
		artist.Scan(&artistname)
		temp := Songschild{Song: artistname + ":Artist"}
		t.MySongChildData = append(t.MySongChildData, temp)

	}
	x := 0
	genrelist := [6]string{"Hip-Hop", "Pop", "Rock", "Country", "Classical", "Jazz"}
	for x < 6 {
		temp := Songschild{Song: genrelist[x] + ":Genre"}
		t.MySongChildData = append(t.MySongChildData, temp)
		x++
	}

	SELECTEDSONGFROMFORM := r.FormValue("songID")
	fmt.Println("IM HERE AHAHAH IM HERE ON POST FOR SONG ID: ", r.FormValue("songID"))
	if r.FormValue("songID") != "" {
		songatt := strings.Split(SELECTEDSONGFROMFORM, "]")
		conv, err := strconv.Atoi(songatt[0])
		getListens := mydb.mydb.QueryRow("SELECT listens from SONG WHERE songID = ?", conv)
		var listenNum int
		getListens.Scan(&listenNum)

		myupdate, err := mydb.mydb.Prepare("UPDATE SONG SET listens = ? WHERE songID = ?")
		defer myupdate.Close()
		myupdateEXE, err := myupdate.Exec(listenNum+1, conv)
		fmt.Println("my execution of this: ", myupdateEXE)

		if err != nil {
			fmt.Println(err)
		}

		var title string
		var path string
		myquery := mydb.mydb.QueryRow("SELECT title,mp3_file from SONG WHERE songID=?", songatt[0])
		myquery.Scan(&title, &path)
		t.SelectedSong = title
		t.SelectedSongPath = path
		if ok {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}
		return
	}

	//THIS WILL SHOW YOU THE SEARCH RESULTS FROM INPUT
	fmt.Println("here")
	searchType := r.FormValue("selectedOptionNAME")
	search := strings.ToLower(r.FormValue("mysearch"))
	fmt.Println(string(searchType))
	fmt.Println("before")
	if strings.Contains(search, ":") {
		search = search[0:strings.Index(search, ":")]
	}
	fmt.Println("THIS IS YOUR SEARCH RESULT: ", search)
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

		var songID int

		for songsID.Next() {
			songsID.Scan(&songID)
			songDetails := mydb.mydb.QueryRow("SELECT title, mp3_file FROM song WHERE songID=?", songID)
			var songTitle string
			var songPathstr string
			if err2 != nil {
				fmt.Println(err2)
			}

			songDetails.Scan(&songTitle, &songPathstr)

			if search == strings.ToLower(songTitle[strings.Index(songTitle, "-")+1:]) {
				fmt.Println("FOUND lsk;dhflkjsdhflkjsd")
				temp := Songschild{SongID: songID, Song: songTitle, SongPath: songPathstr}
				t.MySongChild = append(t.MySongChild, temp)
				fmt.Println(songPathstr)
			}

		}
		if ok {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

	} else if searchType == "Genre" {
		fmt.Println("i'm here now on genre name")
		songs, err2 := mydb.mydb.Query("SELECT songID, title, mp3_file FROM song WHERE genre=?", search)

		if err2 != nil {
			fmt.Println(songs)
		}
		defer songs.Close()
		// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
		// if err2 != nil {
		// 	fmt.Println(songs)
		// }

		var title string
		var mp3file string
		var songIDval int
		for songs.Next() {
			songs.Scan(&songIDval, &title, &mp3file)
			temp := Songschild{SongID: songIDval, Genre: search, Song: title, SongPath: mp3file}
			t.MySongChild = append(t.MySongChild, temp)
		}

		if ok {
			if err := tpl.ExecuteTemplate(w, "search.html", t); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "search_nologin.html", t); err != nil {
				fmt.Println(err)
			}
		}

	} else if searchType == "Artist" {
		fmt.Println("i'm here now")
		artistID, err2 := mydb.mydb.Query("SELECT ArtistID FROM artist")
		if err2 != nil {
			fmt.Println(artistID)
		}
		defer artistID.Close()
		// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
		// if err2 != nil {
		// 	fmt.Println(songs)
		// }

		var artistIDVal int

		for artistID.Next() {
			artistID.Scan(&artistIDVal)
			songDetails, err2 := mydb.mydb.Query("SELECT songID,title, mp3_file FROM song WHERE artist_id=?", artistIDVal)

			var songTitle string
			var songPathstr string
			var songID int
			if err2 != nil {
				fmt.Println(err2)
			}
			for songDetails.Next() {
				songDetails.Scan(&songID, &songTitle, &songPathstr)
				if search == strings.ToLower(songTitle[0:strings.Index(songTitle, "-")]) {
					temp := Songschild{SongID: songID, Song: songTitle, SongPath: songPathstr}
					t.MySongChild = append(t.MySongChild, temp)
					fmt.Println(songPathstr)
				}
			}

		}
		if ok {
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
		artist, err2 := mydb.mydb.Query("SELECT artist_name FROM artist")
		if err2 != nil {
			fmt.Println(artist)
		}
		defer songs.Close()
		defer artist.Close()
		var t Songs

		var songName string

		for songs.Next() {
			songs.Scan(&songName)
			temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:] + ":Song"}
			t.MySongChildData = append(t.MySongChildData, temp2)

		}
		var artistname string
		for artist.Next() {
			artist.Scan(&artistname)
			temp := Songschild{Song: artistname + ":Artist"}
			t.MySongChildData = append(t.MySongChildData, temp)

		}
		x := 0
		genrelist := [6]string{"Hip-Hop", "Pop", "Rock", "Country", "Classical", "Jazz"}
		for x < 6 {
			temp := Songschild{Song: genrelist[x] + ":Genre"}
			t.MySongChildData = append(t.MySongChildData, temp)
			x++
		}
		// newlist := []string{}
		if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
			fmt.Println(err)
		}
		// tpl.ExecuteTemplate(w, "search.html", nil)
		return
	}
	r.ParseForm()
	//datalist for search
	songs, err2 := mydb.mydb.Query("SELECT title FROM song")
	if err2 != nil {
		fmt.Println(songs)
	}
	artist, err2 := mydb.mydb.Query("SELECT artist_name FROM artist")
	if err2 != nil {
		fmt.Println(artist)
	}
	defer songs.Close()
	defer artist.Close()
	var t Songs

	var songName string

	for songs.Next() {
		songs.Scan(&songName)
		temp2 := Songschild{Song: songName[strings.Index(songName, "-")+1:] + ":Song"}
		t.MySongChildData = append(t.MySongChildData, temp2)

	}
	var artistname string
	for artist.Next() {
		artist.Scan(&artistname)
		temp := Songschild{Song: artistname + ":Artist"}
		t.MySongChildData = append(t.MySongChildData, temp)

	}
	x := 0
	genrelist := [6]string{"Hip-Hop", "Pop", "Rock", "Country", "Classical", "Jazz"}
	for x < 6 {
		temp := Songschild{Song: genrelist[x] + ":Genre"}
		t.MySongChildData = append(t.MySongChildData, temp)
		x++
	}

	fmt.Println("here")

	//datalist for search results
	searchType := r.FormValue("selectedOptionNAME")
	search := strings.ToLower(r.FormValue("mysearch"))
	if strings.Contains(search, ":") {
		search = search[0:strings.Index(search, ":")]
	}
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

			var songID int

			for songsID.Next() {
				songsID.Scan(&songID)
				songDetails := mydb.mydb.QueryRow("SELECT title, mp3_file FROM song WHERE songID=?", songID)

				var songTitle string
				var songPathstr string
				if err2 != nil {
					fmt.Println(err2)
				}

				songDetails.Scan(&songTitle, &songPathstr)
				fmt.Println("song title from query ", songTitle)
				fmt.Println("song search from input ", search)
				if search == strings.ToLower(songTitle[strings.Index(songTitle, "-")+1:]) {
					temp := Songschild{SongID: songID, Song: songTitle, SongPath: songPathstr}
					t.MySongChild = append(t.MySongChild, temp)
					fmt.Println(songPathstr)
				}

			}
			if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
				fmt.Println(err)
			}
		} else if searchType == "Genre" {
			fmt.Println("i'm here now on genre name")
			songs, err2 := mydb.mydb.Query("SELECT title, mp3_file FROM song WHERE genre=?", search)

			if err2 != nil {
				fmt.Println(songs)
			}
			defer songs.Close()
			// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
			// if err2 != nil {
			// 	fmt.Println(songs)
			// }

			var title string
			var mp3file string
			for songs.Next() {
				songs.Scan(&title, &mp3file)
				temp := Songschild{Genre: search, Song: title, SongPath: mp3file}
				t.MySongChild = append(t.MySongChild, temp)
			}

			if err := tpl.ExecuteTemplate(w, "searchPlaylist.html", t); err != nil {
				fmt.Println(err)
			}
		} else if searchType == "Artist" {

			fmt.Println("i'm here now")
			artistID, err2 := mydb.mydb.Query("SELECT ArtistID FROM artist")
			if err2 != nil {
				fmt.Println(artistID)
			}
			defer artistID.Close()
			// songPath, err2 := mydb.mydb.mydb.Query("SELECT mp3_file FROM song ")
			// if err2 != nil {
			// 	fmt.Println(songs)
			// }

			var artistIDVal int

			for artistID.Next() {
				artistID.Scan(&artistIDVal)
				songDetails, err2 := mydb.mydb.Query("SELECT songID,title, mp3_file FROM song WHERE artist_id=?", artistIDVal)

				var songTitle string
				var songPathstr string
				var songID int
				if err2 != nil {
					fmt.Println(err2)
				}
				for songDetails.Next() {
					songDetails.Scan(&songID, &songTitle, &songPathstr)
					if search == strings.ToLower(songTitle[0:strings.Index(songTitle, "-")]) {
						temp := Songschild{SongID: songID, Song: songTitle, SongPath: songPathstr}
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
		myquery := mydb.mydb.QueryRow("SELECT title, mp3_file FROM SONG WHERE songID=?", songidconv)
		var title string
		var path string
		myquery.Scan(&title, &path)
		fmt.Println("query for songID after click on results ", title, path)
		insert, err := mydb.mydb.Prepare("INSERT INTO `playlist_song` (`playlistID`, `song`, `songpath`, `songID_pl`) VALUES (?, ?, ?,?);")
		defer insert.Close()
		fmt.Println("THIS IS PLAYLIST VALUE ON PLAYLISTSONGSEARCH SELECTION: ", mydb.playlistID)
		plvalue, err := strconv.Atoi(mydb.playlistID)
		fmt.Println("my playlist ID ahahhahah ", plvalue)
		res, err := insert.Exec(plvalue, title, path, songidconv)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func (mydb *dbstruct) uploadsong(w http.ResponseWriter, r *http.Request) {
	isFound := false
	session, _ := store.Get(r, "session")
	myuserID, _ := session.Values["userID"]
	fmt.Println(myuserID, "userid on open")
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
	defer searchArtist.Close()
	for searchArtist.Next() {
		var t Artist
		searchArtist.Scan(&t.artist)
		if t.artist == r.FormValue("artist_name") {
			isFound = true
		}
	}
	fmt.Println("after search artist next")
	if !isFound {
		insertARTIST, err2 := mydb.mydb.Prepare("INSERT INTO `artist` (`artist_name`) VALUES (?);")
		defer insertARTIST.Close()
		fmt.Println("after prep")
		res2, err2 := insertARTIST.Exec(string(artist_name))
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
	myartistID := mydb.mydb.QueryRow(artistIDQuery)
	fmt.Println("after query artistid")
	artistidNUM := 0
	myartistID.Scan(&artistidNUM)

	fmt.Println("after scan")
	insert, err2 := mydb.mydb.Prepare("INSERT INTO `song` (`release_date`, `title`, `mp3_file`, `UserID`, `listens`, `artist_id`, `genre`) VALUES (?, ?, ?, ?, ?, ?, ?);")
	defer insert.Close()
	fmt.Println("after prep insert song")

	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("no error")
	}
	fmt.Println(myuserID, " this is my user id for insert on song")
	res, err := insert.Exec(time.Now(), string(artist_name+"-"+title), "/songs/"+string(artist_name+"-"+title)+".mp3", myuserID, 1, artistidNUM, string(r.FormValue("genre")))
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
	session, _ := store.Get(r, "session")
	myuserID, ok := session.Values["userID"]
	if r.Method == "GET" {
		if !ok {
			tpl.ExecuteTemplate(w, "home_nologin.html", nil)
			return
		}
		if reflect.ValueOf(myuserID).Interface().(int) == 14 {
			http.Redirect(w, r, "/login.html", 302)
			return
		}

		myquery, err2 := mydb.mydb.Query("SELECT Playlist_ID,playlist_name FROM PLAYLIST WHERE UserID=?", myuserID)

		if err2 != nil {
			fmt.Println(err2)
		}
		defer myquery.Close()
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
			defer myquery.Close()
			var songname string
			var songpath string
			songnameList := ""
			songpathList := ""
			for myquery.Next() {
				//add songs under playlist name and id
				myquery.Scan(&songname, &songpath)
				songnameList = songnameList + songname + ","
				songpathList = songpathList + songpath + ","

			}
			temp := playlistChild{Playlistname: playlistname, PlaylistID: playlistid, Song: songnameList, SongPath: songpathList}
			t.Plchild = append(t.Plchild, temp)
		}
		// top5query, err := mydb.mydb.Query("SELECT song_id, listensTOP FROM TOP_5_SONGS ORDER BY listensTOP DESC")
		// defer top5query.Close()
		// var id int
		// var listens int

		// for top5query.Next() {
		// 	top5query.Scan(&id, &listens)
		// 	myquery := mydb.mydb.QueryRow("SELECT title FROM SONG WHERE songID=?", id)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	var title string
		// 	myquery.Scan(&title)

		// 	temp := Songschild{Song: title, Listens: listens}
		// 	t.Top5 = append(t.Top5, temp)
		// }
		// if err != nil {
		// 	fmt.Println(err)
		// }
		username := mydb.mydb.QueryRow("SELECT username FROM user WHERE UserID=?", myuserID)
		var usernameVal string
		username.Scan(&usernameVal)
		t.Username = usernameVal
		tpl.ExecuteTemplate(w, "home.html", t)
	}
	isPlaylistDelete := r.FormValue("isPlaylistDelete")
	// isSongDelete := r.FormValue("isSongDelete")
	playlistIDVar := r.FormValue("playlistVals")
	fmt.Println(playlistIDVar)
	res := strings.Split(playlistIDVar, "]")
	mydb.playlistID = res[0]
	if isPlaylistDelete == "true" {
		value, err := strconv.Atoi(mydb.playlistID)

		deleteSongsplaylist, err := mydb.mydb.Prepare("DELETE FROM playlist_song WHERE playlistID=?")
		deleteSongsEXE, err := deleteSongsplaylist.Exec(value)
		fmt.Println("delete playlist songs: ", deleteSongsEXE)
		defer deleteSongsplaylist.Close()
		fmt.Println(deleteSongsEXE)

		deleteplaylist, err := mydb.mydb.Prepare("DELETE FROM playlist WHERE Playlist_ID=?")
		deleteEXE, err := deleteplaylist.Exec(value)
		defer deleteplaylist.Close()
		fmt.Println("delete playlist: ", deleteEXE)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/searchPlaylistSong.html", http.StatusSeeOther)

	}
}
func (mydb *dbstruct) createPlaylist(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	myuserID, _ := session.Values["userID"]
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "createplay.html", nil)
		return
	}
	playlistName := r.FormValue("playlistTitle")

	myquery, err := mydb.mydb.Prepare("INSERT INTO `playlist` (`date_created`, `playlist_name`, `UserID`) VALUES (?, ?, ?);")
	defer myquery.Close()
	res, err := myquery.Exec(time.Now(), string(playlistName), myuserID)
	if err != nil {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "createplay.html", "err: "+err.Error())
		return
	}
	print(res)
	tpl.ExecuteTemplate(w, "createplay.html", "Playlist Created")
}
func (mydb *dbstruct) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	http.Redirect(w, r, "/home.html", http.StatusSeeOther)

}
func (mydb *dbstruct) reportsResults(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "sortByDate.html", mydb.Mysongs); err != nil {
		fmt.Println("err")
	}
}
func (mydb *dbstruct) reportsResultsUserCreate(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "sortByUserCreate.html", mydb.UserReport); err != nil {
		fmt.Println("err")
	}
}
func (mydb *dbstruct) reportsResultsPlaylistCreate(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "sortByPlaylistCreate.html", mydb.UserReport); err != nil {
		fmt.Println("err")
	}
}
func (mydb *dbstruct) top5(w http.ResponseWriter, r *http.Request) {
	top5query, err := mydb.mydb.Query("SELECT song_id, listensTOP FROM TOP_5_SONGS ORDER BY listensTOP DESC")
	defer top5query.Close()
	var id int
	var listens int
	var t playlist
	for top5query.Next() {
		top5query.Scan(&id, &listens)
		myquery := mydb.mydb.QueryRow("SELECT title FROM SONG WHERE songID=?", id)
		if err != nil {
			fmt.Println(err)
		}
		var title string
		myquery.Scan(&title)

		temp := Songschild{Song: title, Listens: listens}
		t.Top5 = append(t.Top5, temp)
	}
	if err != nil {
		fmt.Println(err)
	}
	tpl.ExecuteTemplate(w, "top5.html", t)

}
func (mydb *dbstruct) reports(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("here")
		getNames, err := mydb.mydb.Query("SELECT username, UserID FROM user")
		defer getNames.Close()
		if err != nil {
			fmt.Println(err)
		}
		var uName string
		var uID int
		var t userList
		for getNames.Next() {

			getNames.Scan(&uName, &uID)
			temp := userListCHILD{Username: uName, UserID: uID}
			t.Users = append(t.Users, temp)

		}
		if err := tpl.ExecuteTemplate(w, "reports.html", t); err != nil {
			fmt.Println(err)
		}
		return
	}
	r.ParseForm()
	reportMode := r.FormValue("reportMode")
	checkbox := r.FormValue("checkboxvalue")
	userSelect := r.FormValue("userSelect")
	dateCheckbox := r.FormValue("checkboxvaluealltime")
	datefrom := r.FormValue("datefrom") + " 00:00:00"
	dateTo := r.FormValue("dateTo") + " 23:59:59"
	timefrom, err := time.Parse("2006-01-02 15:04:05", datefrom)
	timeto, err := time.Parse("2006-01-02 15:04:05", dateTo)
	fmt.Println("my checkbox for USERS: ", checkbox)
	fmt.Println("my checkbox for DATECHECKBOX: ", dateCheckbox)

	if reportMode == "uploadbydate" {
		var myquery *sql.Rows
		var err2 error
		if checkbox == "Notchecked" {

			fmt.Print("herjhswjhsdfjhzsldkfh")
			conv, err := strconv.Atoi(userSelect)
			myquery, err2 = mydb.mydb.Query("SELECT * FROM song WHERE UserID=?", conv)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			myquery, err2 = mydb.mydb.Query("SELECT * FROM song")
			if err2 != nil {
				fmt.Println(err2)
			}
		}

		defer myquery.Close()
		var songID int
		date := time.Now().Format("2006-01-05 15:04:05")
		var title string
		var mp3_file string
		var userID int
		var listens int
		var artistID int
		var genre string
		var t Songs

		for myquery.Next() {
			myquery.Scan(&songID, &date, &title, &mp3_file, &userID, &listens, &artistID, &genre)
			timeOG, err := time.Parse("2006-01-02 15:04:05", date)
			myqueryName := mydb.mydb.QueryRow("SELECT username FROM user WHERE UserID=?", userID)
			var username string
			myqueryName.Scan(&username)

			if dateCheckbox == "Notchecked" {
				if timeOG.After(timefrom) && timeOG.Before(timeto) {
					temp := Songschild{Uploadedby: username, SongID: songID, Release_date: date, Song: title, UserID: userID, Listens: listens, ArtistID: artistID, Genre: genre}
					t.MySongChild = append(t.MySongChild, temp)
				}
			} else {
				temp := Songschild{Uploadedby: username, SongID: songID, Release_date: date, Song: title, UserID: userID, Listens: listens, ArtistID: artistID, Genre: genre}
				t.MySongChild = append(t.MySongChild, temp)
			}

			if err != nil {
				fmt.Println(err)
			}
		}
		if err != nil {
			fmt.Println(err)
		}

		mydb.Mysongs = t
		http.Redirect(w, r, "/reportsResults.html", http.StatusSeeOther)
	} else if reportMode == "usercreationbydate" {
		fmt.Println("ksjhlkjshfd here on usercreation")
		myquery, err2 := mydb.mydb.Query("SELECT UserID, username, date_registered, name_of_user, email FROM user")

		defer myquery.Close()
		if err2 != nil {
			fmt.Println(err2)
		}
		var userID int
		var username string
		date := time.Now().Format("2006-01-05 15:04:05")
		var fullname string
		var email string
		var t []UserReport
		for myquery.Next() {
			myquery.Scan(&userID, &username, &date, &fullname, &email)
			myquerySongs, err := mydb.mydb.Query("SELECT title FROM song WHERE UserID=?", userID)
			fmt.Println("this is the sought after userid: ", userID)
			defer myquerySongs.Close()
			counter := 0
			for myquerySongs.Next() {
				// myquerySongs.Scan(&titlestring)
				counter++
			}
			timeOG, err := time.Parse("2006-01-02 15:04:05", date)
			if err != nil {
				fmt.Println(err)
			}
			if dateCheckbox == "Notchecked" {
				if timeOG.After(timefrom) && timeOG.Before(timeto) {
					t = append(t, UserReport{Songcounter: strconv.Itoa(counter), Email: email, Userid: userID, Username: username, Date_regist: date, Name_user: fullname})
					// mydb.UserReport = append(mydb.UserReport, UserReport{Userid: userID, Username: username, Date_regist: date, Name_user: fullname})
				}
			} else {
				t = append(t, UserReport{Songcounter: strconv.Itoa(counter), Email: email, Userid: userID, Username: username, Date_regist: date, Name_user: fullname})
				// mydb.UserReport = append(mydb.UserReport, UserReport{Userid: userID, Username: username, Date_regist: date, Name_user: fullname})
			}
		}
		mydb.UserReport = t
		http.Redirect(w, r, "/reportsResultsUserCreate.html", http.StatusSeeOther)
	} else if reportMode == "playlistcreationbydate" {
		fmt.Println("ksjhlkjshfd here")
		myquery, err2 := mydb.mydb.Query("SELECT Playlist_ID, date_created, playlist_name, UserID FROM playlist")

		defer myquery.Close()
		if err2 != nil {
			fmt.Println(err2)
		}
		var userID int
		var playlistid int

		date := time.Now().Format("2006-01-05 15:04:05")
		var playlistname string
		var t []UserReport

		for myquery.Next() {
			myquery.Scan(&playlistid, &date, &playlistname, &userID)
			myqueryName := mydb.mydb.QueryRow("SELECT username FROM user WHERE UserID=?", userID)
			var username string

			myqueryName.Scan(&username)

			plID := strconv.Itoa(playlistid)
			timeOG, err := time.Parse("2006-01-02 15:04:05", date)
			if err != nil {
				fmt.Println(err)
			}
			if dateCheckbox == "Notchecked" {
				if timeOG.After(timefrom) && timeOG.Before(timeto) {
					t = append(t, UserReport{Uploadedby: username, Playlist_id: plID, Userid: userID, Date_regist: date, Playlist_name: playlistname})
				}
			} else {
				t = append(t, UserReport{Uploadedby: username, Playlist_id: plID, Userid: userID, Date_regist: date, Playlist_name: playlistname})
			}
		}
		mydb.UserReport = t
		http.Redirect(w, r, "/reportsResultsPlaylistCreate.html", http.StatusSeeOther)
	}

}
func (mydb *dbstruct) AccountPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("im here")
	session, _ := store.Get(r, "session")
	myuserID, _ := session.Values["userID"]
	// myemail, _ := session.Values["email"]

	type temp struct {
		UserID      int
		Username    string
		Password    string
		Email       string
		Date_regist string
		First       string
		Last        string
	}
	mytest := reflect.ValueOf(myuserID)
	// mytest2 := reflect.ValueOf(myemail)
	myquery := mydb.mydb.QueryRow("SELECT username, email, date_registered, name_of_user FROM USER WHERE UserID=?", myuserID)

	var username string
	date := time.Now().Format("2006-01-02 15:04:05")
	var fullname string
	var email string
	myquery.Scan(&username, &email, &date, &fullname)

	nameSplit := strings.Split(fullname, "_")

	tempstr := temp{UserID: int(mytest.Int()), Username: username, Email: email, Date_regist: date, First: nameSplit[0], Last: nameSplit[1]}
	tpl.ExecuteTemplate(w, "myaccount.html", tempstr)
}
func (mydb *dbstruct) admin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userID, _ := session.Values["userID"]
	mytest := reflect.ValueOf(userID)
	var God HolyGrail
	if mytest.Interface().(int) != 14 {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	ALLSONGS, err := mydb.mydb.Query("SELECT songID, title FROM song")
	if err != nil {
		God.Message = err.Error()
		tpl.ExecuteTemplate(w, "Admin.html", God)
		fmt.Println(err)
		return
	}
	var songID int
	var songName string

	for ALLSONGS.Next() {
		ALLSONGS.Scan(&songID, &songName)
		temp := Songschild{SongID: songID, Song: songName}
		God.MySongChild = append(God.MySongChild, temp)

	}

	ALLUSERS, err := mydb.mydb.Query("SELECT UserID, username FROM user")
	if err != nil {
		God.Message = err.Error()
		tpl.ExecuteTemplate(w, "Admin.html", God)
		fmt.Println(err)
		return
	}
	var userIDVAL int
	var username string
	for ALLUSERS.Next() {
		ALLUSERS.Scan(&userIDVAL, &username)
		if userIDVAL != 14 {
			temp := UserReport{Userid: userIDVAL, Username: username}
			God.MyUserReports = append(God.MyUserReports, temp)
		}

	}

	// ALLPLAYLISTS, err := mydb.mydb.Query("SELECT Playlist_ID, playlist_name, UserID FROM playlist")
	ALLPLAYLISTS, err := mydb.mydb.Query("SELECT playlist.Playlist_ID, playlist.playlist_name, playlist.UserID, user.username FROM `3380-project`.playlist,`3380-project`.user WHERE user.UserID=playlist.UserID")

	if err != nil {
		God.Message = err.Error()
		tpl.ExecuteTemplate(w, "Admin.html", God)
		fmt.Println(err)
		return
	}
	var playlistID int
	var playlistname string
	var userIDpl int

	var usernamePL string
	for ALLPLAYLISTS.Next() {
		ALLPLAYLISTS.Scan(&playlistID, &playlistname, &userIDpl, &usernamePL)

		temp := playlistChild{Username: usernamePL, PlaylistID: playlistID, Playlistname: playlistname}
		God.MyPlaylistChilds = append(God.MyPlaylistChilds, temp)
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "Admin.html", God)

	}

	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("on POST")
		selectedType := r.FormValue("selectedType")
		fmt.Println(selectedType)
		if selectedType == "songRemove" {
			songID := r.FormValue("songID")
			fmt.Println(songID)
			// songName := r.FormValue("songName")

			deleteSongPlaylist, err := mydb.mydb.Prepare("DELETE FROM playlist_song WHERE songID_pl=?")
			songIDtoINT, err := strconv.Atoi(songID)
			path := mydb.mydb.QueryRow("SELECT mp3_file FROM song WHERE songID=?", songIDtoINT)
			var mp3path string

			path.Scan(&mp3path)
			fmt.Println(mp3path[1:])
			e := os.Remove(mp3path[1:])

			if e != nil {
				fmt.Println(e)
			}
			if err != nil {
				fmt.Println(err)
			}

			top5song := mydb.mydb.QueryRow("SELECT song_id FROM top_5_songs WHERE song_id=?", songIDtoINT)
			var songIDtop5 int
			top5error := top5song.Scan(&songIDtop5)

			if top5error == nil { //SONGID IS IN TOP 5
				fmt.Println("song is in top 5 ")
				fmt.Println("im here on delete")
				lastplacenewVal := mydb.mydb.QueryRow("SELECT songID,listens FROM song ORDER BY listens DESC LIMIT 5,1;")
				var songIDlast int
				var listensLast int
				lastplacenewVal.Scan(&songIDlast, &listensLast)
				fmt.Println("sdfsdf" + strconv.Itoa(songIDlast))
				insertSongTop5, err := mydb.mydb.Prepare("INSERT INTO top_5_songs (song_id,listensTOP) VALUES (?,?);")
				_, myerr := insertSongTop5.Exec(songIDlast, listensLast)
				if myerr != nil {
					God.Message = myerr.Error()
					fmt.Println(myerr)
					tpl.ExecuteTemplate(w, "Admin.html", God)
					return

				}

				deleteSongTop5, err := mydb.mydb.Prepare("DELETE FROM top_5_songs WHERE song_id=?")
				if err != nil {
					fmt.Println(err)
				}
				_, errdel := deleteSongTop5.Exec(&songIDtop5)
				fmt.Println(errdel)

			} else {
				fmt.Println("not on top 5")
				fmt.Println(top5error)
			}
			deleteSong, err := mydb.mydb.Prepare("DELETE FROM song WHERE songID=?")

			if err != nil {
				God.Message = err.Error()
				fmt.Println(err)
				tpl.ExecuteTemplate(w, "Admin.html", God)

				return
			}
			_, errpl := deleteSongPlaylist.Exec(songIDtoINT)
			_, errsong := deleteSong.Exec(songIDtoINT)
			fmt.Println(errpl)
			fmt.Println(errsong)
			fmt.Println("song removed")
			God.Message = "song removed"

		} else if selectedType == "userRemove" {
			userID := r.FormValue("userID")
			// username := r.FormValue("username")
			fmt.Println("im on user remove")
			songIDtoINT, err := strconv.Atoi(userID)
			getPlaylistSongs, err := mydb.mydb.Query("SELECT DISTINCT playlist_song.PlaylistID FROM `3380-project`.playlist, `3380-project`.playlist_song WHERE playlist.UserID=? AND playlist.Playlist_ID=playlist_song.PlaylistID;", songIDtoINT)

			var playlistIDplaylistsong int
			for getPlaylistSongs.Next() {
				getPlaylistSongs.Scan(&playlistIDplaylistsong)
				deleteUserPlaylist, err := mydb.mydb.Prepare("DELETE FROM playlist_song WHERE playlistID=?")
				if err != nil {
					fmt.Println(err)
				}
				_, err2 := deleteUserPlaylist.Exec(playlistIDplaylistsong)
				fmt.Print(err2)
			}
			deleteUserPlaylist, err := mydb.mydb.Prepare("DELETE FROM playlist WHERE UserID=?")
			_, myerr1 := deleteUserPlaylist.Exec(songIDtoINT)
			if myerr1 != nil {
				fmt.Println(myerr1)
			}
			deleteUserSongs, err := mydb.mydb.Prepare("DELETE FROM song WHERE UserID=?")

			_, myerr3 := deleteUserSongs.Exec(songIDtoINT)
			if myerr1 != nil {
				fmt.Println(myerr1)
			}
			if myerr3 != nil {
				fmt.Println(myerr3)
			}
			deleteUser, err := mydb.mydb.Prepare("DELETE FROM user WHERE UserID=?")

			if err != nil {
				God.Message = err.Error()
				tpl.ExecuteTemplate(w, "Admin.html", God)
				fmt.Println(err)
				return
			}

			_, myerr := deleteUser.Exec(songIDtoINT)
			if myerr != nil {
				fmt.Println(myerr)
			}
			God.Message = "User removed"

		} else if selectedType == "playlistRemove" {
			playlistID := r.FormValue("playlistID")
			// playlistName := r.FormValue("playlistName")

			deleteSongPlaylist, err := mydb.mydb.Prepare("DELETE FROM playlist_song WHERE playlistID=?")
			playlistIDtoINT, err := strconv.Atoi(playlistID)
			_, errDel := deleteSongPlaylist.Exec(playlistIDtoINT)
			fmt.Println(errDel)
			deletePlaylist, err := mydb.mydb.Prepare("DELETE FROM playlist WHERE Playlist_ID=?")
			_, errdelpl := deletePlaylist.Exec(playlistIDtoINT)
			fmt.Println(errdelpl)
			if err != nil {
				God.Message = err.Error()
				tpl.ExecuteTemplate(w, "Admin.html", God)
				fmt.Println(err)
				return
			}

			God.Message = "Playlist removed"

		}

		tpl.ExecuteTemplate(w, "Admin.html", God)

	}

}
func (mydb *dbstruct) editP(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userID, _ := session.Values["userID"]
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "editP.html", nil)
		return
	}
	fmt.Println("Im at POST of editP")
	username := r.FormValue("username")
	email := r.FormValue("email")
	var myquery *sql.Stmt
	var err error

	queryUName, err := mydb.mydb.Query("SELECT username FROM USER")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer queryUName.Close()

	for queryUName.Next() {
		var t User
		queryUName.Scan(&t.username)
		if t.username == username {
			tpl.ExecuteTemplate(w, "editP.html", "This username already exists")
			return
		}
	}

	query, err := mydb.mydb.Query("SELECT email FROM USER")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer query.Close()

	for query.Next() {
		var t User
		query.Scan(&t.email)
		if t.email == email {
			tpl.ExecuteTemplate(w, "editP.html", "This Account already exists")
			return
		}
	}

	if err != nil {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "editP.html", err)
		return
	}
	if username == "" {
		myquery, err = mydb.mydb.Prepare("UPDATE user SET username=username, email=? WHERE UserID=?")
		defer myquery.Close()

		res, err := myquery.Exec(email, userID)
		fmt.Println(res, err)
	} else if email == "" {
		myquery, err = mydb.mydb.Prepare("UPDATE user SET username=?, email=email WHERE UserID=?")
		fmt.Println("IM ON EMAIL EMPTY :))")
		defer myquery.Close()

		res, err := myquery.Exec(username, userID)
		fmt.Println(res, err)
	} else if email == "" && username == "" {
		tpl.ExecuteTemplate(w, "editP.html", "Both Fields cannot be empty")
		return
	} else {
		myquery, err = mydb.mydb.Prepare("UPDATE user SET username=?, email=? WHERE UserID=?")
		fmt.Println("IM ON EMAIL EMPTY :))")
		defer myquery.Close()

		res, err := myquery.Exec(username, email, userID)
		fmt.Println(res, err)
	}

	tpl.ExecuteTemplate(w, "editP.html", "Successful")
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

	query, err := mydb.mydb.Query("SELECT email FROM USER")
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

	queryUName, err := mydb.mydb.Query("SELECT username FROM USER")
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

	insert, err2 := mydb.mydb.Prepare("INSERT INTO `user` (`username`, `password`, `email`, `date_registered`, `name_of_user`, `access_level`) VALUES (?, ?, ?, ?, ?, '1');")
	defer insert.Close()
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
	tpl.ExecuteTemplate(w, "signup.html", "Account Creation Successful")
	// defer insert.Close()

	fmt.Println("sdkjfhsdf here")

}
func (mydb *dbstruct) Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		val, ok := session.Values["userID"]
		mytest := reflect.ValueOf(val)
		if !ok || mytest.Interface().(int) == 14 {
			http.Redirect(w, r, "/login.html", 302)
			return
		}
		// ServeHTTP calls f(w, r)
		// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
		HandlerFunc.ServeHTTP(w, r)
	}
}
func AuthAdmin(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		val, ok := session.Values["userID"]
		mytest := reflect.ValueOf(val)
		if !ok || mytest.Interface().(int) != 14 {
			http.Redirect(w, r, "/login.html", 302)
			return
		}
		// ServeHTTP calls f(w, r)
		// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
		HandlerFunc.ServeHTTP(w, r)
	}
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

	mydb := dbstruct{mydb: myDB()}
	// test := webHandler{mu: mux}
	http.HandleFunc("/login.html", mydb.login)
	http.HandleFunc("/signup", mydb.addAccountSignUp)

	http.HandleFunc("/logout", mydb.logout)
	http.HandleFunc("/", mydb.home)
	http.HandleFunc("/upload_song", mydb.Auth(mydb.uploadsong))
	http.HandleFunc("/search.html", mydb.Auth(mydb.searchList))
	http.HandleFunc("/search_nologin.html", mydb.searchList)

	http.HandleFunc("/searchPlaylistSong.html", mydb.Auth(mydb.searchListPlaylist))
	http.HandleFunc("/forgotpassword.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/forgotpassword.html")
	})
	http.HandleFunc("/Admin.html", AuthAdmin(mydb.admin))
	http.HandleFunc("/editP.html", mydb.Auth(mydb.editP))
	http.HandleFunc("/AO.html", mydb.Auth(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/AO.html")
	}))
	http.HandleFunc("/sortByDate.html", mydb.Auth(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/sortByDate.html")
	}))
	http.HandleFunc("/reports", AuthAdmin(mydb.reports))
	http.HandleFunc("/top5.html", mydb.Auth(mydb.top5))
	http.HandleFunc("/reportsResults.html", AuthAdmin(mydb.reportsResults))
	http.HandleFunc("/reportsResultsUserCreate.html", AuthAdmin(mydb.reportsResultsUserCreate))
	http.HandleFunc("/reportsResultsPlaylistCreate.html", AuthAdmin(mydb.reportsResultsPlaylistCreate))

	http.HandleFunc("/changepass.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/changepass.html")
	})
	http.HandleFunc("/myaccount.html", mydb.Auth(mydb.AccountPage))
	http.HandleFunc("/createplay.html", mydb.Auth(mydb.createPlaylist))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images"))))
	http.Handle("/songs/", http.StripPrefix("/songs/", http.FileServer(http.Dir("./songs"))))
	listenAddr := fmt.Sprintf(":%d", args.port)

	log.Fatal(http.ListenAndServe("127.0.0.1"+listenAddr, context.ClearHandler(http.DefaultServeMux)))

}

/*To run website server, type localhost:8085 on browser, it will run a local version of the website.
No one else will have access to it but on your machine

**test**
*/
