package main

import (
        "encoding/json"
        "log"
        "net/http"
        "io/ioutil"
        "github.com/cgundogan/KDTree/kdtree"
)

type JsonMsg struct {
        Vectors []kdtree.Vector
        FindNearestTo kdtree.Vector
}

func handlerKDTree(w http.ResponseWriter, req *http.Request) {
        var msg JsonMsg
        
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                panic(err)
        }
        log.Println(string(body))
        if err := json.Unmarshal(body, &msg); err != nil {
                panic(err)
        }
        
        dataSet := kdtree.NewDataSet(msg.Vectors, len(msg.Vectors[0]), 0)
        kdTree := kdtree.NewKDTreeByDataSet(dataSet)
        foundVector := kdTree.FindNearest(msg.FindNearestTo)

        js, err := json.Marshal(foundVector)
        if err != nil {
                panic(err)
        }

        log.Println(msg)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write(js)
}

func main() {
        log.Println("starting server..")
        http.HandleFunc("/kdtree/ajax", handlerKDTree)
        http.Handle("/", http.FileServer(http.Dir("assets")))
        log.Fatal(http.ListenAndServe(":8080", nil))
}
