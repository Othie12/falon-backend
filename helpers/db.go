package helpers

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Initializing connection to database
func Init() (*sql.DB, error) {
	var err error

	//Openning a connection to mysql database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/falon")
	if err != nil {
		return nil, err
	}

	//checking connection to verify it it worked
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//return the database connection if it worked
	return db, nil
}

// function to put files into the file server
// after it's been resolved from the http post request
func Upload(filepath string, file multipart.File) (bool, error) {
	outfile, err := os.Create(filepath) //create a file in the desired filepath
	if err != nil {
		return false, err
	}
	defer outfile.Close()

	//copy the file that's been uploaded into the file that was created on the server
	_, err = io.Copy(outfile, file)
	if err != nil {
		return false, err
	}

	//return true if the file was uploaded successfully
	return true, nil
}

func Upload_improved(wantedFilepath string, file multipart.File, header *multipart.FileHeader) (string, error) {
	//generate filename and save file to disk
	fileExt := filepath.Ext(header.Filename)
	var item string
	if fileExt == ".mp3" {
		item = "beat"
	} else {
		item = "pic"
	}
	fileName := fmt.Sprint(item, time.Now().UnixNano(), fileExt)

	f, err := os.OpenFile(wantedFilepath+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)

	if err != nil {
		return "", err
	}
	return fileName, nil
}
