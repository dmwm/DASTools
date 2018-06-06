package main

// DASTools: mongoimport upload JS files into MongoDB
// Copyright (c) 2018 - Valentin Kuznetsov <vkuznet AT gmail dot com>

import (
    "bufio"
	"flag"
	"fmt"
    "encoding/json"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var port int
	flag.IntVar(&port, "port", 8230, "MongoDB port")
	var host string
	flag.StringVar(&host, "host", "", "MongoDB hostname")
	var db string
	flag.StringVar(&db, "db", "", "MongoDB db name")
	var collection string
	flag.StringVar(&collection, "collection", "", "MongoDB collection name")
	var file string
	flag.StringVar(&file, "file", "", "file to upload")
	flag.Parse()
	if version {
		fmt.Println(info())
		return
	}
	uri := fmt.Sprintf("mongodb://%s:%d", host, port)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer session.Close()

    c := session.DB(db).C(collection)
    f, err := os.Open(file)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        var data map[string]interface{}
        err := json.Unmarshal([]byte(line), &data)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        c.Insert(data)
    }
}
