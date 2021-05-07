package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var dbFileName = "/code/ScoreKeeper.db"

func main() {
	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	/*
		https://www.ietf.org/rfc/rfc3986.txt
		https://tools.ietf.org/html/rfc6570
		The body is the data.
		The query string is the query.
			If you looking to return a 404 then the query data should be part of the URL
				/game/123
			If you are looking to return an empty list then use a parameter
				/scores?game_id=123
			You have to ask yourself: is it an identifier or a filter
	*/
	// Handling the endpoints
	// User
	router.HandleFunc("/api/users", GetUsers).Methods("GET")
	router.HandleFunc("/api/user/{id}", GetUser).Methods("GET")
	router.HandleFunc("/api/user", AddUser).Methods("POST")

	// Game
	router.HandleFunc("/api/games", GetGames).Methods("GET")
	router.HandleFunc("/api/game/{id}", GetGame).Methods("GET")
	router.HandleFunc("/api/game", AddGame).Methods("POST")

	// Scores
	router.HandleFunc("/api/scores", GetScores).Methods("GET")
	router.HandleFunc("/api/score", SetScore).Methods("POST")
	router.HandleFunc("/api/score", UpdateScore).Methods("PUT")

	// Running the Server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: "No user 'id' provided"})
	} else {
		db := DB{Filename: dbFileName}
		userId, convertErr := strconv.Atoi(vars["id"])
		if convertErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while attempting to convert user id: " + convertErr.Error()})
		} else {
			user, err := db.GetUserById(int32(userId))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while looking up a user: " + err.Error()})
			} else if (user == User{}) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUsers")
	w.Header().Set("Content-Type", "application/json")

	db := DB{Filename: dbFileName}
	defer db.closeConnection()

	users := db.GetUsers()
	if users == nil {
		users = []User{}
	}

	json.NewEncoder(w).Encode(users)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("AddUser called")
	w.Header().Set("Content-Type", "application/json")

	var userInfo User
	json.NewDecoder(r.Body).Decode(&userInfo)

	// At least a first name needs to be provided.
	if userInfo.FirstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: "firstname can not be blank"})
	} else {
		db := DB{Filename: dbFileName}

		userId, err := db.CreateUser(userInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while adding a user: " + err.Error()})
		} else {
			userInfo.Id = userId
			userInfo.Active = 1

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(userInfo)
		}
	}
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: "No game 'id' provided"})
	} else {
		db := DB{Filename: dbFileName}
		gameId, convertErr := strconv.Atoi(vars["id"])
		if convertErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while attempting to convert game id: " + convertErr.Error()})
		} else {
			game, err := db.GetGameById(int32(gameId))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while looking up a game: " + err.Error()})
			} else if (game == Game{}) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				json.NewEncoder(w).Encode(game)
			}
		}
	}
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	log.Println("GetGames")
	w.Header().Set("Content-Type", "application/json")

	db := DB{Filename: dbFileName}
	defer db.closeConnection()

	games := db.GetGames()
	if games == nil {
		games = []Game{}
	}

	json.NewEncoder(w).Encode(games)
}

func AddGame(w http.ResponseWriter, r *http.Request) {
	log.Println("AddGame called")
	w.Header().Set("Content-Type", "application/json")

	var gameInfo Game
	json.NewDecoder(r.Body).Decode(&gameInfo)

	// At least a first name needs to be provided.
	if gameInfo.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: "Game name can not be blank"})
	} else {
		db := DB{Filename: dbFileName}

		gameId, err := db.CreateGame(gameInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while adding a game: " + err.Error()})
		} else {
			gameInfo.Id = gameId

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(gameInfo)
		}
	}
}

func GetScores(w http.ResponseWriter, r *http.Request) {
	log.Println("GetGames")
	w.Header().Set("Content-Type", "application/json")

	db := DB{Filename: dbFileName}
	defer db.closeConnection()

	scores := db.GetScores()
	if scores == nil {
		scores = []Score{}
	}

	json.NewEncoder(w).Encode(scores)
}

func SetScore(w http.ResponseWriter, r *http.Request) {
	log.Println("SetScore called")
	w.Header().Set("Content-Type", "application/json")

	var scoreInfo Score
	json.NewDecoder(r.Body).Decode(&scoreInfo)

	ok, err := scoreInfo.Validate()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: err.Error()})
	} else {
		db := DB{Filename: dbFileName}

		scoreId, err := db.InsertScore(scoreInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while adding new score: " + err.Error()})
		} else {
			scoreInfo.Id = scoreId

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(scoreInfo)
		}
	}
}

func UpdateScore(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateScore called")
	w.Header().Set("Content-Type", "application/json")

	var scoreInfo Score
	json.NewDecoder(r.Body).Decode(&scoreInfo)

	ok, err := scoreInfo.Validate()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Error: err.Error()})
	} else {
		db := DB{Filename: dbFileName}

		_, err := db.UpdateScore(scoreInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMsg{Error: "Error while adding new score: " + err.Error()})
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(scoreInfo)
		}
	}
}

type ErrorMsg struct {
	Error string `json:"error"`
}

type User struct {
	Id        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    int32  `json:"active"`
}

type Game struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Score struct {
	Id         int32  `json:"id"`
	UserId     int32  `json:"user_id"`
	GameId     int32  `json:"game_id"`
	DatePlayed string `json:"date_played"`
	Score      int32  `json:"score"`
}

type DB struct {
	Filename   string
	connection *sql.DB
}

func (s Score) Validate() (bool, error) {
	var valid bool = false
	var retError error
	var errMsgs []string

	if s.DatePlayed == "" {
		dt := time.Now()
		s.DatePlayed = dt.Format("2021-05-11")
	} else {
		correctFormat, err := regexp.Match(`^\d{4}-\d{2}-\d{2}$`, []byte(s.DatePlayed))
		if err != nil {
			errMsgs = append(errMsgs, "Unrecoverable error while validating date played information: "+err.Error())
		}
		if !correctFormat {
			errMsgs = append(errMsgs, "Date Played needs to be in the correct format: YYYY-MM-DD")
		}
	}

	if s.GameId == 0 {
		errMsgs = append(errMsgs, "Game ID needs to be provided")
	}

	if s.UserId == 0 {
		errMsgs = append(errMsgs, "User ID needs to be provided")
	}

	if len(errMsgs) == 0 {
		valid = true
	} else {
		retError = errors.New(strings.Join(errMsgs, "\n"))
	}

	return valid, retError
}

func (db *DB) initDB() {
	log.Println("initializing the DB connection")
	connection, err := sql.Open("sqlite3", db.Filename)
	if err != nil {
		panic(err)
	}
	if connection == nil {
		panic("db nil")
	}
	db.connection = connection
}

func (db *DB) getConnection() *sql.DB {
	log.Println("Getting DB Connection")
	if db.connection == nil {
		db.initDB()
	}

	return db.connection
}

func (db *DB) closeConnection() {
	if db.connection != nil {
		db.connection.Close()
	}
}

func (db DB) CreateUser(user User) (int32, error) {
	var userId int64

	sql := `
		INSERT INTO user
		(first_name, last_name, active)
		values (?, ?, 1)
	`
	stmt, err := db.getConnection().Prepare(sql)
	if err == nil {
		result, err2 := stmt.Exec(user.FirstName, user.LastName)
		if err2 != nil {
			err = err2
		} else {
			userId, err = result.LastInsertId()
		}
	}

	return int32(userId), err
}

func (db DB) GetUsers() (users []User) {
	log.Println("DB GetUser")
	sql := `SELECT id, first_name, last_name, active FROM user`

	dbh := db.getConnection()
	rows, err := dbh.Query(sql)

	if err != nil {
		log.Println("Error while attempting to get all users:")
		panic(err)
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Active)
		users = append(users, user)
	}

	return users
}

func (db DB) GetUserById(id int32) (User, error) {
	log.Println("DB GetUserById")
	sql := `SELECT id, first_name, last_name, active FROM user WHERE id = ?`

	dbh := db.getConnection()
	rows, err := dbh.Query(sql, id)

	if err != nil {
		log.Println("Error while attempting to get all users:")
		panic(err)
	}

	var user User
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Active)
	}

	return user, err
}

func (db DB) CreateGame(game Game) (int32, error) {
	var gameId int64

	sql := `
		INSERT INTO game
		(name)
		values (?)
	`
	stmt, err := db.getConnection().Prepare(sql)
	if err == nil {
		result, err2 := stmt.Exec(game.Name)
		if err2 != nil {
			err = err2
		} else {
			gameId, err = result.LastInsertId()
		}
	}

	return int32(gameId), err
}

func (db DB) GetGames() (games []Game) {
	log.Println("DB GetUser")
	sql := `SELECT id, name FROM game`

	dbh := db.getConnection()
	rows, err := dbh.Query(sql)

	if err != nil {
		log.Println("Error while attempting to get all games:")
		panic(err)
	}

	for rows.Next() {
		game := Game{}
		err = rows.Scan(&game.Id, &game.Name)
		games = append(games, game)
	}

	return games
}

func (db DB) GetGameById(id int32) (Game, error) {
	log.Println("DB GetGameById")
	sql := `SELECT id, name FROM game WHERE id = ?`

	dbh := db.getConnection()
	rows, err := dbh.Query(sql, id)

	if err != nil {
		log.Println("Error while attempting to get specific game:")
		panic(err)
	}

	var game Game
	if rows.Next() {
		err = rows.Scan(&game.Id, &game.Name)
	}

	return game, err
}

func (db DB) GetScores() (scores []Score) {
	log.Println("DB GetScores")
	sql := `SELECT id, user_id, game_id, date_played, score FROM score`

	dbh := db.getConnection()
	rows, err := dbh.Query(sql)

	if err != nil {
		log.Println("Error while attempting to get all scores:")
		panic(err)
	}

	for rows.Next() {
		score := Score{}
		err = rows.Scan(&score.Id, &score.UserId, &score.GameId, &score.DatePlayed, &score.Score)
		scores = append(scores, score)
	}

	return scores
}

func (db DB) InsertScore(score Score) (int32, error) {
	var scoreId int64

	sql := `
		INSERT INTO score
		(user_id, game_id, date_played, score)
		values (?, ?, ?, ?)
	`
	stmt, err := db.getConnection().Prepare(sql)
	if err == nil {
		result, err2 := stmt.Exec(score.UserId, score.GameId, score.DatePlayed, score.Score)
		if err2 != nil {
			err = err2
		} else {
			scoreId, err = result.LastInsertId()
		}
	}

	return int32(scoreId), err
}

func (db DB) UpdateScore(score Score) (bool, error) {
	var ok bool

	sql := `
		UPDATE score
		SET user_id = ?, game_id = ?, date_played = ?, score = ?
		WHERE id = ?
	`
	stmt, err := db.getConnection().Prepare(sql)
	if err == nil {
		_, err2 := stmt.Exec(score.UserId, score.GameId, score.DatePlayed, score.Score, score.Id)
		if err2 != nil {
			err = err2
		} else {
			ok = true
		}
	}

	return ok, err
}
