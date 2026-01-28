package main

import (
 //   "io"
    "fmt"
	"strings"
	"strconv"
	"path"
    "net/http"
	"github.com/amadeon/megregator/internal/model"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func getStrSlice(slice []string, index int) (string, bool) {
	if index >= 0 && index < len(slice) {
		return slice[index], true
	}
	return "", false
}

func defHandler(res http.ResponseWriter, req *http.Request){
	http.NotFound(res, req)
}

func upHandler(res http.ResponseWriter, req *http.Request){
    res.Header().Set("Content-Type", "text/html")

	//if req.Method != http.MethodPost{
	//	http.NotFound(res, req)
	//	return
	//}

	d := strings.Split(path.Clean(req.URL.Path), "/")

	t, exists := getStrSlice(d, 0)
	typ := models.MetrType(t)
	if !exists || typ != models.TCounter && typ != models.TGauge{
		http.Error(res, "bad value", http.StatusBadRequest)
		return
	}
	name, exists := getStrSlice(d, 1)
	if !exists{
		http.NotFound(res, req)
		return
	}
	val, exists := getStrSlice(d, 2)
	if !exists || !isNumeric(val){
		http.Error(res, "bad value", http.StatusBadRequest)
		return
	} 


    fmt.Fprintf(res, "%v %v %v\n\n",typ, name, val)
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", defHandler)

	mux.Handle("/update/", http.StripPrefix(`/update/`, http.HandlerFunc(upHandler)))

	err := http.ListenAndServe(`:8080`, mux)
    if( err!= nil ){
		panic(err)
    }

}