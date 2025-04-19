package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"

	//	"strconv"
	//	"strings"
	"time"
)

/*
func CorsMiddleWare(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Set to "true" if credentials are required
n(w, r)
	}
}
*/

func GetUserFromContext(r *http.Request) *User {
	const userKey string = "user"
	userid, ok := r.Context().Value(userKey).(string)

	if !ok {
		fmt.Println("user id not found in context")
		return nil
	}

	user := User{
		Id: userid,
	}

	return &user
}

func CheckSession(shouldexist bool, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Figure out environment variables for this
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL")) // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		db := GetDb()
		sessionId, err := r.Cookie("session_id")
		if err != nil {
			if err == http.ErrNoCookie {
				if !shouldexist {
					h(w, r)
					return
				}
				fmt.Println("this is the redirect")
				http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusSeeOther)
			}
			fmt.Fprintln(w, err.Error())
			return
		}
		if shouldexist {
			var mysession Session
			mysession.Id = sessionId.Value
			dSession := db.QueryRow("SELECT expires_at FROM session WHERE id = ?", mysession.Id)
			dSerr := dSession.Scan(&mysession.ExpiresAt)

			if dSerr != nil {
				fmt.Fprintln(w, dSerr.Error())
			}

			// Delete session if expired, update active session if halfway dead
			if mysession.ExpiresAt < time.Now().Unix() {
				// Session is expired
				_, derr := db.Exec("DELETE FROM session WHERE id = ?", mysession.Id)
				if derr != nil {
					fmt.Fprintln(w, derr.Error())
				}
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else if mysession.ExpiresAt-time.Now().Unix() < (7 * 24 * 60 * 60) {
				newTime := time.Now().Unix() + (15 * 24 * 60 * 60)
				_, derr := db.Exec("UPDATE session SET expires_at = ? WHERE id = ?", newTime, mysession.Id)
				if derr != nil {
					fmt.Fprintln(w, derr.Error())
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    mysession.Id,
					Path:     "/",
					MaxAge:   15 * 24 * 60 * 60,
					HttpOnly: true,
				})
			}

			const userKey string = "user"
			var myuser User
			dUser := db.QueryRow("SELECT user_id FROM session WHERE id = ?", mysession.Id)
			dUerr := dUser.Scan(&myuser.Id)
			if dUerr != nil {
				fmt.Fprintln(w, dUerr.Error())
				return
			}

			ctx := context.WithValue(r.Context(), userKey, myuser.Id)
			h(w, r.WithContext(ctx))
			return
		}
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"/app", http.StatusSeeOther)
	}
}

func CreateRoutes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	mux.HandleFunc("/auth/google/callback", CheckSession(false, CallbackHandle))
	mux.HandleFunc("/logout/google", LogoutHandle)
	mux.HandleFunc("/auth/google", CheckSession(false, AuthHandle))
	mux.HandleFunc("GET /app", CheckSession(true, appGet))
	mux.HandleFunc("POST /app", CheckSession(true, appPost))
	mux.HandleFunc("POST /workouts", CheckSession(true, workoutsPost))

	return mux
}

func CallbackHandle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	tuser, gerr := gothic.CompleteUserAuth(w, r)
	if gerr != nil {
		fmt.Fprintln(w, gerr)
		return
	}

	// Check if user exists in db
	var myuser User
	db := GetDb()
	dUser := db.QueryRow("SELECT id, google_email FROM user WHERE google_email = ?", tuser.Email)
	dberr := dUser.Scan(&myuser.Id, &myuser.GoogleEmail) // See if query can fit into a User if not make an error

	if dberr != nil {
		if dberr == sql.ErrNoRows {
			// User does not exist error
			userId := uuid.New()
			_, uerr := db.Exec("INSERT INTO user (id, google_email) VALUES (?, ?)", userId, tuser.Email)

			if uerr != nil {
				fmt.Fprintln(w, uerr)
				return
			}
			myuser.Id = userId.String()

			// Insert default stuff in
			_, eerr := db.Exec("INSERT INTO exercise (id, user_id, name) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?)", uuid.New(), myuser.Id, "Bench", uuid.New(), myuser.Id, "Squat", uuid.New(), myuser.Id, "Deadlift")
			if eerr != nil {
				fmt.Fprintln(w, eerr)
				return
			}
		} else {
			// DB just failed or smthn
			fmt.Fprintln(w, dberr.Error())
			return
		}
	}

	// Create session
	fmt.Println(myuser.Id)
	sessionId := uuid.New()
	expiresAt := time.Now().Unix() + (15 * 24 * 60 * 60) // Make session expire in 15 days
	fmt.Println(expiresAt)
	_, serr := db.Exec("INSERT INTO session (id, user_id, expires_at) VALUES (?, ?, ?)", sessionId, myuser.Id, expiresAt)
	if serr != nil {
		fmt.Fprintln(w, serr)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		Path:     "/",
		MaxAge:   15 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   false, // Change to true in production
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"/app", http.StatusTemporaryRedirect)
}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	sessionId, _ := r.Cookie("session_id")

	db := GetDb()
	db.Exec("DELETE FROM session WHERE id = ?", sessionId)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Change in prod
	})

	gothic.Logout(w, r)
	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusTemporaryRedirect)
}

func AuthHandle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"/app", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func appGet(w http.ResponseWriter, r *http.Request) {
	myuser := GetUserFromContext(r)
	if myuser == nil {
		fmt.Println("error context was nil")
		return
	}

	db := GetDb()

	var name sql.NullString

	dUser := db.QueryRow("SELECT name FROM user WHERE id = ?", myuser.Id)

	duNerr := dUser.Scan(&name)
	if duNerr != nil {
		fmt.Fprintln(w, duNerr.Error())
		return
	}

	myuser.Name = name.String

	var exercises []Exercise
	dExcersises, deError := db.Query("SELECT id, name FROM exercise WHERE user_id = ?", myuser.Id)
	if deError != nil {
		fmt.Fprintln(w, deError.Error())
		return
	}

	for dExcersises.Next() {
		var ex Exercise
		err := dExcersises.Scan(&ex.Id, &ex.Name)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		exercises = append(exercises, ex)
	}

	w.Header().Set("Content-Type", "application/json")
	if !name.Valid {
		data := Data{
			NameNotValid: true,
		}

		jsonbytes, jerr := json.Marshal(data)
		if jerr != nil {
			http.Error(w, jerr.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonbytes)
		return
	}

	data := Data{
		NameNotValid: false,
		Name:         myuser.Name,
		Exercises:    exercises,
	}

	jsonbytes, jerr := json.Marshal(data)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonbytes)
}

type Body struct {
	Create     bool   `json:"create"`
	Delete     bool   `json:"delete"`
	CreateName bool   `json:"createName"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	Weight     int    `json:"weight"`
	Sets       int    `json:"sets"`
	Reps       int    `json:"reps"`
}

func appPost(w http.ResponseWriter, r *http.Request) {
	myuser := GetUserFromContext(r)
	if myuser == nil {
		fmt.Println("error context was nil")
		return
	}
	decoder := json.NewDecoder(r.Body)

	var body Body
	berr := decoder.Decode(&body)
	if berr != nil {
		http.Error(w, berr.Error(), http.StatusBadRequest)
	}

	validate := validator.New()
	db := GetDb()

	if body.CreateName {
		verr := validate.Var(body.Name, "ascii")
		if verr != nil {
			http.Error(w, "name not accepted", http.StatusNotAcceptable)
		}
		_, derr := db.Exec("UPDATE user SET name = ? WHERE user_id = ?", body.Name, myuser.Id)
		if derr != nil {
			fmt.Fprintln(w, derr.Error())
			return
		}
		return
	}

	if body.Create {
		verr := validate.Var(body.Name, "ascii")
		if verr != nil {
			http.Error(w, "name not accepted", http.StatusNotAcceptable)
			return
		}

		_, derr := db.Exec("INSERT INTO exercise (id, user_id, name) VALUES (?, ?, ?)", uuid.New(), myuser.Id, body.Name)
		if derr != nil {
			fmt.Fprintln(w, derr.Error())
			return
		}

		return
	}

	if body.Delete {
		_, derr := db.Exec("DELETE FROM exercise WHERE id = ?", body.Id)
		if derr != nil {
			fmt.Fprintln(w, derr.Error())
			return
		}
		return
	}

	verr := validate.StructPartial(body, "Body.Id Body.Weight Body.Sets Body.Reps")
	if verr != nil {
		http.Error(w, "not accepted", http.StatusNotAcceptable)
		return
	}

	_, derr := db.Exec("INSERT INTO workout (id, user_id, exercise_id, weight, sets, reps, time) VALUES (?, ?, ?, ?, ?, ?, ?)", uuid.New(), myuser.Id, body.Id, body.Weight, body.Sets, body.Reps, time.Now().Unix())
	if derr != nil {
		fmt.Fprintln(w, derr.Error())
		return
	}
}

type WorkoutResponse struct {
	Workouts []Workout `json:"workouts"`
}

func workoutsPost(w http.ResponseWriter, r *http.Request) {
	myuser := GetUserFromContext(r)
	if myuser == nil {
		fmt.Println("error context was nil")
		return
	}
	db := GetDb()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the body", http.StatusBadRequest)
		return
	}

	exerciseId := string(body)
	defer r.Body.Close()

	var rightUser int
	dRightUser := db.QueryRow("SELECT count(1) FROM exercise WHERE id = ? AND user_id = ?", exerciseId, myuser.Id)
	dRightUser.Scan(&rightUser)
	if rightUser != 1 {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var workouts []Workout
	dWorkouts, dwError := db.Query("SELECT weight, reps, sets, time FROM workout WHERE exercise_id = ?", exerciseId)
	if dwError != nil {
		fmt.Fprintln(w, dwError.Error())
		return
	}

	for dWorkouts.Next() {
		var work Workout
		err := dWorkouts.Scan(&work.Weight, &work.Reps, &work.Sets, &work.Time)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		workTime := time.Unix(work.Time, 0)
		workString := workTime.Month().String() + " " + strconv.Itoa(workTime.Day()) + ", " + strconv.Itoa(workTime.Year())

		work.TimeString = workString
		workouts = append(workouts, work)
	}

	w.Header().Set("Content-Type", "application/json")

	data := WorkoutResponse{
		Workouts: workouts,
	}

	jsonbytes, jerr := json.Marshal(data)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonbytes)
}
