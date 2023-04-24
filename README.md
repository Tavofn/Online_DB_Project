﻿<a name="br1"></a>ReadMe

So the files we are submitting are the SQL dump file named “[team3_dump.sql](https://github.com/Tavofn/Online_DB_Project/blob/main/sql/dump/team3_dump.sql)[ ](https://github.com/Tavofn/Online_DB_Project/blob/main/sql/dump/team3_dump.sql)“, and our
folders that contain our front-end and back-end code.

The folders we are submitting include the following:
<pre>
● Server folder 
   ○ main.go - Contains all of our backend code written in Go.

● SQL folder 
  ○ Dump - Folder which contains our populate SQL dump file
   ■ team3_dump.sql
   ○ Mini_world.sql 
   ○ Testsongs.csv - folder of song data we used for testing purposes

● Web folder 
   ○ This folder contains all of our frontend code and all the images we used for our website
</pre>
Our project requires the installation of golang. If you prefer the see the hosted website, all you will need is the
link to our website(https://team3.coogsmusic.com/) and some account login information,
which will be provided in another document in our submitted folder. Although, if you
want to run our code locally there are some steps that are required:

1\. Install https://go.dev/doc/install

2\. After you have all of our code in your developer environment, open up a terminal
 and enter this command “Go run ./server/main.go”

3\. After that wait a couple of seconds and you will see a message in the terminal
 that says “Successful Connection to Database!”.

4\. Then, open up your browser and type “<http://localhost:8086>”

5\. From there you can navigate our site locally
