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
	//     "github.com/ghodss/yaml"
)

type Params map[string]string
type DASMap struct {
	System    string   `yaml:"system" json:"system"`
	Format    string   `yaml:"format" json:"format"`
	Instances []string `yaml:"instances" json:"instances"`
	Url       string   `yaml:"urn" json:"urn"`
	Urn       string   `yaml:"url" json:"url"`
	Lookup    string   `yaml:"lookup" json:"lookup"`
	Expire    int      `yaml:"expire" json:"expire"`
	Params    Params   `yaml:"params" json:"params"`
	Das_map   []Params `yaml:"das_map" json:"das_map"`
    Type      string   `yaml:"type" json:"type"`
    Hash      string   `yaml:"hash" json:"hash"`
    TimeStamp int64    `yaml:"ts" json:"ts"`
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
        r, e := json.Marshal(d)
        if e != nil {
            log.Fatal(e)
        }
		d.Hash = fmt.Sprintf("%x", md5.Sum(r))
        d.TimeStamp = time.Now().Unix()
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
