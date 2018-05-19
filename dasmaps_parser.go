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

// String method provides string representation of DASMap
func (d *DASMap) String() string {
	return recString(d)
}

// Notation represents notation DAS record
type Notation struct {
	Api        string `yaml:"api" json:"api"`
	Api_output string `yaml:"api_output" json:"api_output"`
	Rec_key    string `yaml:"rec_key" json:"rec_key"`
}

// String method provides string representation of DASMap
func (n *Notation) String() string {
	return recString(n)
}

// Notations represents collection of notation DAS records
type Notations struct {
	Hash      string     `yaml:"hash" json:"hash"`
	Notations []Notation `yaml:"notations" json:"notations"`
	TimeStamp int64      `yaml:"ts" json:"ts"`
	Type      string     `yaml:"type" json:"type"`
}

// String method provides string representation of DASMap
func (n *Notations) String() string {
	return recString(n)
}

// InputValue represents input_values DAS record
type InputValue struct {
	Input             string `yaml:"input" json:"input"`
	Jsonpath_selector string `yaml:"jsonpath_selector" json:"jsonpath_selector"`
	Test              string `yaml:"test" json:"test"`
	Url               string `yaml:"url" json:"url"`
}

// String method provides string representation of DASMap
func (n *InputValue) String() string {
	return recString(n)
}

// InputValues represents collection of input_values DAS records
type InputValues struct {
	Hash        string       `yaml:"hash" json:"hash"`
	InputValues []InputValue `yaml:"input_values" json:"input_values"`
	TimeStamp   int64        `yaml:"ts" json:"ts"`
	Type        string       `yaml:"type" json:"type"`
}

// String method provides string representation of DASMap
func (n *InputValues) String() string {
	return recString(n)
}

// String method provides string representation of DASMap
func recString(v interface{}) string {
	r, e := json.Marshal(v)
	if e != nil {
		log.Fatal(e)
	}
	return string(r)
}

// helper function to parse input yaml file and produce DAS map records
func dasmaps(input string) {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	var gRec DASMap
	for i, r := range strings.Split(string(data), "---") {
		if strings.Contains(r, "urn") { // das record
			n := DASMap{}
			err = yaml.Unmarshal([]byte(r), &n)
			if err != nil {
				log.Fatalf("record: %v, %v", r, err)
			}
			if i == 0 {
				gRec = n
				continue
			}
			n.System = gRec.System
			n.Format = gRec.Format
			n.Instances = gRec.Instances
			n.Type = "service"
			n.TimeStamp = time.Now().Unix()
			r, e := json.Marshal(n)
			if e != nil {
				log.Fatal(e)
			}
			n.Hash = fmt.Sprintf("%x", md5.Sum(r))
			fmt.Println(n.String())
		} else if strings.Contains(r, "notations") {
			n := Notations{}
			err = yaml.Unmarshal([]byte(r), &n)
			if err != nil {
				log.Fatalf("record: %v, %v", r, err)
			}
			n.Type = "notation"
			n.TimeStamp = time.Now().Unix()
			r, e := json.Marshal(n)
			if e != nil {
				log.Fatal(e)
			}
			n.Hash = fmt.Sprintf("%x", md5.Sum(r))
			fmt.Println(n.String())
		} else if strings.Contains(r, "input_values") {
			n := InputValues{}
			err = yaml.Unmarshal([]byte(r), &n)
			if err != nil {
				log.Fatalf("record: %v, %v", r, err)
			}
			n.Type = "input_values"
			n.TimeStamp = time.Now().Unix()
			r, e := json.Marshal(n)
			if e != nil {
				log.Fatal(e)
			}
			n.Hash = fmt.Sprintf("%x", md5.Sum(r))
			fmt.Println(n.String())
		}
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
