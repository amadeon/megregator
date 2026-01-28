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
	if len(d) < 3{
		http.Error(res, "bad arguments list", http.StatusBadRequest)
		return
	}
	typ, name, val := models.MetrType(d[0]), d[1], d[2] 
	if typ != models.TCounter && typ != models.TGauge{
		http.NotFound(res, req)
		return
	}
	if !isNumeric(val){
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