package main

// DASTools: mongostatus returns string with basic MongoDB stats
// Copyright (c) 2018 - Valentin Kuznetsov <vkuznet AT gmail dot com>

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8230, "MongoDB port")
	flag.Parse()
	uri := fmt.Sprintf("mongodb://localhost:%d", port)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer session.Close()

	r := bson.M{}
	if err := session.DB("admin").Run(bson.D{{"serverStatus", 1}}, &r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proc, _ := r["process"]
	pid, _ := r["pid"]
	ver, _ := r["version"]
	host, _ := r["host"]
	c, _ := r["connections"]
	con := c.(bson.M)
	cons := con["totalCreated"]
	s, _ := r["ok"]
	status := fmt.Sprintf("%v", s)
	if s == 1. {
		status = "ok"
	}
	fmt.Printf("Process:%s, PID:%d, Version:%s, Host:%s, Connections:%d, Status:%s\n", proc, pid, ver, host, cons, status)
}
