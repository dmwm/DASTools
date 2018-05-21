package main

// DASTools: das_cleanup clean up DAS cache/merge collections
// Copyright (c) 2018 - Valentin Kuznetsov <vkuznet AT gmail dot com>

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func cleanup(port, delta, verbose int) {
	uri := fmt.Sprintf("mongodb://localhost:%d", port)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer session.Close()
	cond := bson.M{"$lt": time.Now().Unix() - int64(delta)}
	spec := bson.M{"das.expire": cond}
	cols := []string{"cache", "merge"}
	for _, col := range cols {
		c := session.DB("das").C(col)
		n, err := c.RemoveAll(spec)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if verbose > 0 {
			fmt.Printf("remove %d docs from %s\n", n.Removed, col)
		}
	}
}

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
	var delta int
	flag.IntVar(&delta, "delta", 3600, "delta TTL value")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	if version {
		info()
		return
	}
	cleanup(port, delta, verbose)
}
