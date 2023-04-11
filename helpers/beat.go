package helpers

import (
	"database/sql"
	"net/http"
)

// The beats class...
type Beat struct {
	Id, Bpm, Owner_id                                                                   uint32
	Name, Genre, Date_uploaded, Duration, Description, Artwork_filepath, Audio_filepath string
	Price                                                                               float64
}

// register a file to the database
func (b *Beat) Register_to_database(db *sql.DB) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO beat(producer_id, name, description, genre, duration, bpm, price, audio_file, artwork) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Owner_id, b.Name, b.Description, b.Genre, b.Duration, b.Bpm, b.Price, b.Audio_filepath, b.Artwork_filepath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// method to handle beat uploads from client. But it still needs me to read all of the required data and store it
// into the beat struct
// need to also turn the other functions that need not to be exported into lowercase first letters
func (b *Beat) Upload_beat_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(30 << 20)                // Limit file size to 30 MB
	file, handler, err := r.FormFile("audioFile") // Extract the uploaded file from the request
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest) // Return an error if the file is invalid
		return
	}
	defer file.Close()

	// Save file to disk
	beat_filepath := "/beats/" + handler.Filename // Define the file path on the server

	_, err = Upload(beat_filepath, file)
	if err != nil {
		http.Error(w, "Failed to upload ", http.StatusBadRequest)
		return
	}

	//Initialize connection to db. Using the Init function I made earlier in db.go
	db, err := Init()
	if err != nil {
		http.Error(w, "Failed to connect db ", http.StatusBadRequest)
		return
	}

	//taking the info to the db using a method I created earlier
	_, err = b.Register_to_database(db)
	if err != nil {
		http.Error(w, "Failed to save information to database ", http.StatusBadRequest)
	}

}
