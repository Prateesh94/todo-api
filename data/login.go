package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo-api/encrypt"

	"github.com/gorilla/mux"
)

type Login struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"description"`
	email string
}

var jwtkey = []byte("pi-is-infinite")

func verifyBody(a Login) error {
	if a.Name == "" || a.Email == "" || a.Password == "" {
		return errors.New("invalid data:- found empty fields")
	}
	return nil
}

func AddNewUserEndpoint(w http.ResponseWriter, r *http.Request) {
	body := json.NewDecoder(r.Body)
	var dat Login
	er := body.Decode(&dat)
	if er != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid JSON")
		return
	}
	err := verifyBody(dat)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(w, err)
		return
	}
	hsh := encrypt.Crypt(dat.Password)
	val := fmt.Sprintf("%x", hsh)
	err = addUser(dat.Name, dat.Email, val)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, err)
		return
	}
}
func LoginEndpoint(w http.ResponseWriter, r *http.Request) {
	var usr Login
	body := json.NewDecoder(r.Body)
	body.Decode(&usr)
	usr.Name = "default"
	er := verifyBody(usr)
	if er != nil {
		fmt.Fprintln(w, er)
		return
	}
	c := fmt.Sprintf("%x", encrypt.Crypt(usr.Password))
	d := verifyUser(usr.Email, c)
	if d {
		tok, _ := jwtTokenCreator(usr.Email)
		mp := make(map[string]string)
		mp["token"] = tok
		w.Header().Set("Authorization", "Bearer "+tok)
		json.NewEncoder(w).Encode(mp)

	} else {
		msg := make(map[string]string)
		msg["message"] = "User not found"
		json.NewEncoder(w).Encode(msg)
		return
	}
}

func AddTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	a := r.Header.Get("Authorization")
	email, tok, er := authhandler(a)
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := make(map[string]string)
		msg["message"] = "Unauthorized"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tok)

	b := json.NewDecoder(r.Body)
	var dat Todo
	var got Todo
	b.Decode(&dat)
	got, er = addTodo(dat.Title, dat.Desc, email)
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := make(map[string]string)
		msg["message"] = "Unauthorized"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(got)
}

func UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	a := r.Header.Get("Authorization")
	email, tok, er := authhandler(a)
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := make(map[string]string)
		msg["message"] = "Unauthorized"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tok)
	var dat Todo
	k := mux.Vars(r)
	id := k["id"]
	if id == "" {
		fmt.Fprintln(w, "Invalid Id")
		return
	}
	b := json.NewDecoder(r.Body)
	b.Decode(&dat)
	got, er := updateTodo(dat.Title, dat.Desc, id, email)
	if er != nil {
		w.WriteHeader(http.StatusForbidden)
		msg := make(map[string]string)
		msg["message"] = "Forbidden"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(got)
}

func DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	a := r.Header.Get("Authorization")
	email, tok, er := authhandler(a)
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := make(map[string]string)
		msg["message"] = "Unauthorized"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tok)
	k := mux.Vars(r)
	id := k["id"]
	if id == "" {
		fmt.Fprintln(w, "Invalid Id")
		return
	}
	er = deleteTodo(id, email)
	if er != nil {
		w.WriteHeader(http.StatusForbidden)
		msg := make(map[string]string)
		msg["message"] = "Forbidden"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusOK)
}

func FetchtodoEndpoint(w http.ResponseWriter, r *http.Request) {
	a := r.Header.Get("Authorization")
	email, tok, er := authhandler(a)
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := make(map[string]string)
		msg["message"] = "Unauthorized"
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tok)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 5
	}
	offset := (page - 1) * limit
	pgdat, er := fetchtodo(email, limit, offset)
	if er != nil {
		w.WriteHeader(http.StatusNotFound)
		msg := make(map[string]string)
		msg["message"] = "No records found"
		fmt.Fprintf(w, `{"data":%v,"page":%d,"limit":%d}`, msg["message"], page, limit)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	showdat, _ := json.MarshalIndent(pgdat, "", " ")
	fmt.Fprintf(w, `{"data":%s,"page":%d,"limit":%d,"total":%d}`, showdat, page, limit, len(pgdat))
}
