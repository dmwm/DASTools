package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// NOTE: The json package always orders keys when marshalling. Specifically:
// - Maps have their keys sorted lexicographically
// - Structs keys are marshalled in the order defined in the struct

// Params keeps DAS map parameters
type Params map[string]string

// DASMap represent generic DAS map, we use particular order to perform serialization
type DASMap struct {
	Das_map   []Params `yaml:"das_map" json:"das_map"`
	Expire    int      `yaml:"expire" json:"expire"`
	Format    string   `yaml:"format" json:"format"`
	Hash      string   `yaml:"hash" json:"hash"`
	Instances []string `yaml:"instances" json:"instances"`
	Lookup    string   `yaml:"lookup" json:"lookup"`
	Params    Params   `yaml:"params" json:"params"`
	System    string   `yaml:"system" json:"system"`
	TimeStamp int64    `yaml:"ts" json:"ts"`
	Type      string   `yaml:"type" json:"type"`
	Url       string   `yaml:"url" json:"url"`
	Urn       string   `yaml:"urn" json:"urn"`
}

func (d *DASMap) String() string {
	r, e := json.Marshal(d)
	if e != nil {
		log.Fatal(e)
	}
	return string(r)
}

func dasmaps(input string) {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	var gRec DASMap
	for i, r := range strings.Split(string(data), "---") {
		d := DASMap{}
		err = yaml.Unmarshal([]byte(r), &d)
		if err != nil {
			log.Fatalf("record: %v, %v", r, err)
		}
		if i == 0 {
			gRec = d
			continue
		}
		d.System = gRec.System
		d.Format = gRec.Format
		d.Instances = gRec.Instances
		d.Type = "service"
		d.TimeStamp = time.Now().Unix()
		r, e := json.Marshal(d)
		if e != nil {
			log.Fatal(e)
		}
		d.Hash = fmt.Sprintf("%x", md5.Sum(r))
		fmt.Println(d.String())
	}
}

func main() {
	var input string
	flag.StringVar(&input, "input", "", "yam das map file")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	dasmaps(input)
}
