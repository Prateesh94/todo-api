package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-api/data"
	"todo-api/encrypt"

	"github.com/gorilla/mux"
)

var jsonusers = [][]byte{
	[]byte(`{"":"test1",
	"emil":"testt1@case.com",
	"paword":"testcase1"
	}`),
	[]byte(`{"name":"test1",
"email":"testt1@case.com",
"password":"testcase1"
}`),
	[]byte(`{"name":"",
"email":"test2@case.com",
"password":"testcase2"
}`),
	[]byte(`{"name":"test3",
"email":"",
"password":"testcase3"
}`),
	[]byte(`{"name":"test4 Case",
"email":"test4@case.com",
"password":""
}`),
	[]byte(`{"name":"test5",
"email":"test5@case.com",
"password":"testcase5"
}`),
	[]byte(`{"name":"test1",
"email":"testt1@case.com",
"password":"testcase1"
}`), []byte(`{"name":"Name1",
"email":"name1@case.com",
"password":"name1"
}`)}

var jsontodo = [][]byte{
	[]byte(`{"name":"test1",
"email":"testt1@case.com",
"password":"testcase1"
}`),
	[]byte(`{"id":"test1",
"title":"Test_title_1",
"description":"test_description_1"
}`),
	[]byte(`{"title":"Test_title_2",
"description":"test_description_2"
}`),
	[]byte(`{"description":"Test_title_3",
"title":"test_description_3"
}`),
	[]byte(`{"":"Test_title_4",
"description":"test_description_4"
}`),
	[]byte(`{"title":"Test_title_5",
"description":"test_description_5"
}`), []byte(`{"title":"Update1",
"description":"Update1"
}`)}
var token string
var tok []string

func requrl(a, b string, c []byte, f map[string]string) *http.Request {
	r := httptest.NewRequest(a, b, bytes.NewBuffer(c))
	return mux.SetURLVars(r, f)
}
func TestRegister(t *testing.T) {
	for _, k := range jsonusers {
		req, er := http.NewRequest("POST", "/register", bytes.NewBuffer(k))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if er != nil {
			fmt.Println(er)
		}
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(data.AddNewUserEndpoint)
		handle.ServeHTTP(rr, req)
		fmt.Println(rr.Body)
		fmt.Println(rr.Code)
	}

}

func TestLogin(t *testing.T) {
	for _, k := range jsonusers {
		req, er := http.NewRequest("POST", "/login", bytes.NewBuffer(k))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if er != nil {
			fmt.Println(er)
		}
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(data.LoginEndpoint)
		handle.ServeHTTP(rr, req)
		token = rr.Header().Get("Authorization")
		tok = append(tok, token)
		fmt.Println(rr.Body)
		fmt.Println(rr.Code)
	}
}

func TestAddTodo(t *testing.T) {
	for _, j := range tok {
		for _, k := range jsontodo {
			req, er := http.NewRequest("POST", "/todos", bytes.NewBuffer(k))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			req.Header.Set("Authorization", j)
			if er != nil {
				fmt.Println(er)
			}
			rr := httptest.NewRecorder()
			handle := http.HandlerFunc(data.AddTodoEndpoint)
			handle.ServeHTTP(rr, req)
			fmt.Println(rr.Code)
		}
	}
}

func TestUpdateTodo(t *testing.T) {

	for _, j := range tok {
		for _, k := range jsontodo {

			for ids := 1; ids < 10; ids++ {
				vars := make(map[string]string)
				vars["id"] = fmt.Sprintf("%d", ids)
				req := requrl("PUT", "/todos/{id}", k, vars)
				req.Header.Set("Authorization", j)
				rr := httptest.NewRecorder()
				handle := http.HandlerFunc(data.UpdateEndpoint)
				handle.ServeHTTP(rr, req)
				fmt.Println(rr.Code)
			}
		}
		vars := make(map[string]string)
		vars["id"] = "error"
		req := requrl("PUT", "/todos/{id}", nil, vars)
		req.Header.Set("Authorization", j)
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(data.UpdateEndpoint)
		handle.ServeHTTP(rr, req)
		fmt.Println(rr.Code)
	}
}

func TestDeleteEndpoint(t *testing.T) {
	for _, j := range tok {
		for _, k := range jsontodo {

			for ids := 10; ids < 15; ids++ {
				vars := make(map[string]string)
				vars["id"] = fmt.Sprintf("%d", ids)
				req := requrl("DELETE", "/todos/{id}", k, vars)
				req.Header.Set("Authorization", j)
				rr := httptest.NewRecorder()
				handle := http.HandlerFunc(data.DeleteEndpoint)
				handle.ServeHTTP(rr, req)
				fmt.Println(rr.Code)
			}
		}
		vars := make(map[string]string)
		vars["id"] = "error"
		req := requrl("DELETE", "/todos/{id}", nil, vars)
		req.Header.Set("Authorization", j)
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(data.UpdateEndpoint)
		handle.ServeHTTP(rr, req)
		fmt.Println(rr.Code)
	}
}

func TestFetchtodo(t *testing.T) {

	for _, j := range tok {
		req := httptest.NewRequest("GET", "/todos?page=1&limit=5", nil)
		req.Header.Set("Authorization", j)
		rr := httptest.NewRecorder()
		han := http.HandlerFunc(data.FetchtodoEndpoint)
		han.ServeHTTP(rr, req)
		fmt.Println(rr.Code)
	}
}

func TestLimit(t *testing.T) {
	url := "/register"
	met := "POST"
	fmt.Println("into limiter test")
	for i := 0; i < 55; i++ {
		req, _ := http.NewRequest(met, url, bytes.NewBuffer(jsonusers[4]))
		rr := httptest.NewRecorder()
		handle := http.HandlerFunc(data.AddNewUserEndpoint)
		jan := encrypt.Limitmid(handle)
		jan.ServeHTTP(rr, req)
		fmt.Println(rr.Body)
	}
}
func TestRefreshtoken(t *testing.T) {
	req := httptest.NewRequest("GET", "/todos", nil)
	req.Header.Set("Authorization", tok[1])
	han := http.HandlerFunc(data.FetchtodoEndpoint)
	rr := httptest.NewRecorder()
	han.ServeHTTP(rr, req)
	val := rr.Header().Get("Authorization")
	for val == tok[1] {
		time.Sleep(time.Minute * 3)
		fmt.Println("MAKING REQUEST AGAIN")
		han.ServeHTTP(rr, req)
		val = rr.Header().Get("Authorization")
	}

	fmt.Println(tok[1], "-old \n new- \n", rr.Header().Get("Authorization"))
}
