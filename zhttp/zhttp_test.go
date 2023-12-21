package zhttp_test

import (
	"encoding/json"
	"fmt"
	"github.com/uaxe/infra/zhttp"
	"net/http"
	"net/http/httptest"
)

type Ret struct {
	ContentType string `header:"Content-Type"`
	XRequestID  string `header:"X-Request-ID"`
	Name        string `json:"name,omitempty"`
	Age         int    `json:"age,omitempty"`
}

func ExampleParse() {
	r, _ := http.NewRequest("GET", "/", nil)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", "1")
		rr := Ret{Name: "zkep", Age: 18}
		raw, _ := json.Marshal(rr)
		w.Write(raw)
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	ret := rec.Result()
	var v Ret
	err := zhttp.DefaultRespParse.Parse(ret, &v)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", v.Name)
	fmt.Printf("%#v\n", v.Age)
	fmt.Printf("%#v\n", v.ContentType)
	fmt.Printf("%#v\n", v.XRequestID)
	// Output:
	// <nil>
	// "zkep"
	// 18
	// "application/json"
	// "1"
}
