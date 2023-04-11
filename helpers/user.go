package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// declaring a session store variable
var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

// first declair some variables since we expect to get null values
var id sql.NullString
var fname, lname, email, uname, profile_pic sql.NullString

// the user struct
type User struct {
	Id                                                    string `json:"id"`
	Fname, Lname, Uname, Email, Pwd, Profile_pic_filepath string
}

func (u *User) Register_to_database(db *sql.DB) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO user(fname, lname, uname, email, pwd, profile_pic) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	worked, err := stmt.Exec(u.Fname, u.Lname, u.Uname, u.Email, u.Pwd, u.Profile_pic_filepath)
	if err != nil {
		return false, err
	}
	fmt.Println(worked.LastInsertId())
	return true, nil
}

// function that gets called whenever user fills the login form
func (u *User) Login_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error reading json file"+err.Error(), http.StatusBadRequest)
		log.Println("Error reading json file" + err.Error())
		return
	}
	email, ok := data["email"].(string)
	if !ok {
		log.Println("invalid email")
	}
	pwd, ok := data["pwd"].(string)
	if !ok {
		log.Println("invalid pwd")
	}
	log.Println("email and name: " + email + ", " + pwd)
	db, err := Init()
	if err != nil {
		fmt.Fprint(w, "failed to connect db"+err.Error())
		return
	}
	passed, errorMessage := u.Login(db, email, pwd)
	if !passed {
		http.Error(w, errorMessage, http.StatusBadRequest)
		log.Println(errorMessage)
	} else {
		session, _ := store.Get(r, "user-session")
		session.Values["authenticated"] = true
		session.Values["id"] = u.Id
		log.Println("user id: ", u.Id)
		session.Save(r, w)
		log.Println("Session id: ", session.ID)
		log.Println("session user id: ", session.Values["id"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m := make(map[string]string)
		m["id"] = u.Id
		log.Println(m)
		log.Println("m['id']", m["id"])
		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			log.Println("Error making json file", err.Error())
			http.Error(w, "Error making json file"+err.Error(), http.StatusBadRequest)
		}

	}

}

func (u *User) Login(db *sql.DB, uname string, pwd string) (bool, string) {
	//Database query
	stmt, err := db.Prepare("SELECT * FROM user WHERE email = ? AND pwd = ?")
	if err != nil {
		return false, "Failed to prepare sql stmt: " + err.Error()
	}
	defer stmt.Close()

	//query the db using prepared statemtent
	err = stmt.QueryRow(uname, pwd).Scan(&id, &fname, &lname, &uname, &email, &pwd, &profile_pic)
	if err != nil {
		return false, err.Error()
	}

	if id.Valid {
		u.Id = string(id.String)
	}
	if fname.Valid {
		u.Fname = string(fname.String)
	}
	if lname.Valid {
		u.Lname = string(lname.String)
	}
	if email.Valid {
		u.Email = string(email.String)
	}
	if profile_pic.Valid {
		u.Profile_pic_filepath = string(profile_pic.String)
	}

	fmt.Println("it worked")
	fmt.Printf("id: %d\n Names: %s %s\n uname: %s\n", u.Id, u.Fname, u.Lname, u.Uname)
	return true, ""
}

func (u *User) Logout_handler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Values["authenticated"] = false
	session.Values["id"] = ""
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}

func (u *User) GetUser(db *sql.DB, id string) (User, error) {
	stmt, err := db.Prepare("SELECT id, fname, lname, uname, email, profile_pic FROM user WHERE id = ?")
	if err != nil {
		return *u, err
	}
	defer stmt.Close()
	/*
		//query the db using prepared statemtent
		rows, err := stmt.Query(id).Scan()
		if err != using prepared statemtent
		rows, err := stmt.Query(id).Scan()
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		//Iterate through records
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.id, &user.fname, &user.lname, &user.uname, &user.email, &user.pwd, &user.profile_pic); err != nil {
				panic(err.Error())
			}

			fmt.Println("it worked")
			fmt.Printf("id: %d\n Names: %s %s\n uname: %s\n", user.id, user.fname.String, user.lname.String, user.uname.String)
	*/
	err = stmt.QueryRow(id).Scan(&id, &fname, &lname, &uname, &email, &profile_pic)
	if err != nil {
		return *u, err
	}

	var user User

	user.Id = id

	if fname.Valid {
		user.Fname = string(fname.String)
	}
	if lname.Valid {
		user.Lname = string(lname.String)
	}
	if uname.Valid {
		user.Uname = string(uname.String)
	}
	if email.Valid {
		user.Email = string(email.String)
	}
	if profile_pic.Valid {
		user.Profile_pic_filepath = string(profile_pic.String)
	}

	return user, nil
}

// get user from an http request
func (u *User) Register_user_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(30 << 20)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//first get the form values and assign them to variables
	u.Fname = r.FormValue("fname")
	u.Lname = r.FormValue("lname")
	u.Email = r.FormValue("email")
	u.Pwd = r.FormValue("pwd")
	//u.Contact = r.FormValue("contact")
	u.Uname = r.FormValue("uname")

	//handle files that have been uploaded
	ffile, headr, err := r.FormFile("profilePic")
	if err != nil {
		log.Println("error reading file" + err.Error())
		fmt.Fprintf(w, err.Error())
	}

	if ffile != nil {
		//generate filename and save file to disk
		u.Profile_pic_filepath, err = Upload_improved("profilepics/", ffile, headr)
		if err != nil {
			log.Println("Error saving file into filesystem: ", err.Error())
		}
	}

	db, err := Init()
	if err != nil {
		http.Error(w, "Failed to connect db: ", http.StatusBadRequest)
		fmt.Fprintf(w, "failed to connect db: ", err)
		log.Println("Failed to connect db: ", http.StatusBadRequest, err)
		return
	}

	_, err = u.Register_to_database(db)
	if err != nil {
		http.Error(w, "Failed to save information to database ", http.StatusBadRequest)
		log.Println("Failed to save to db: ", http.StatusBadRequest, err)
		fmt.Fprintf(w, "failed to save db: ", err)
		return
	}
	fmt.Fprintf(w, "Registered successfuly")
	log.Println("registered")
	return
}
