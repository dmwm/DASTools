package main

// DASTools: dasmaps_validator validates DAS maps
// Copyright (c) 2018 - Valentin Kuznetsov <vkuznet AT gmail dot com>

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// validate input dasmap file
/*
{"hash": "84b55197205fc1ef1fcdb943a2487862", "format": "JSON", "url": "https://cmsweb.cern.ch/sitedb/data/prod/roles", "urn": "roles", "ts": 1525096018.0, "system": "sitedb2", "das_map": [{"rec_key": "user.role", "api_arg": "match", "das_key": "role"}], "services": "", "expire": 3600, "lookup": "role", "wild_card": "*", "params": {"match": "optional"}, "type": "service"}
{"arecord": {"count": 6, "type": "service", "system": "sitedb2"}}
{"input_values": [{"url": "https://cmsweb.cern.ch/sitedb/data/prod/site-names", "input": "site.name", "test": "T1*", "jsonpath_selector": "$.result[*][2]"}], "type": "input_values", "system": "sitedb2", "hash": "ed47f30560fda174df1d2f38073a8d02"}
{"notations": [{"rec_key": "site.name", "api": "", "api_output": "site.cms_name"}, {"rec_key": "name", "api": "group_responsibilities", "api_output": "user_group"}, {"rec_key": "name", "api": "", "api_output": "alias"}], "hash": "4b861b22321625e9f9f5b2e25b0c3762", "type": "notation", "ts": 1525096019.0, "system": "sitedb2"}
{"arecord": {"count": 1, "type": "notation", "system": "sitedb2"}}
{"verification_token": "7915820eaee71521c9bb6d7345731666", "type": "verification_token"}
*/
// var (
//     dasMapKeys := []string{"rec_key", "api_arg", "das_key"}
//     dasRecKeys := []string{"hash", "format", "url", "urn", "ts", "system", "das_maps", "services", "expire", "lookup", "wild_card", "params", "type"}
//     notationKeys := []string{"notations", "hash", "type", "ts", "system"}
//     notMapKeys := []string{"rec_key", "api", "api_outptu"}
//     arecordKeys := []string{"count", "type", "system"}
//     verKeys := []string{"verification_token", "type"}
// )

type Record map[string]interface{}

// NOTE: The json package always orders keys when marshalling. Specifically:
// - Maps have their keys sorted lexicographically
// - Structs keys are marshalled in the order defined in the struct
func checkHash(rec Record) bool {
	keys := mapKeys(rec)
	if inList("hash", keys) {
		h := rec["hash"]
		// reset hash value to calculate new md5 checksum
		rec["hash"] = ""
		data, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
		}
		rh := fmt.Sprintf("%x", md5.Sum(data))
		if rh != h {
			fmt.Println(string(data))
			fmt.Println("")
			fmt.Println(" record type", recordType(rec))
			fmt.Println(" record hash", h)
			fmt.Println("computed md5", rh)
			log.Fatal("Invalid hash")
		}
	}
	return true
}

// inList helper function to check item in a list
func inList(a string, list []string) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// mapKeys helper function to return keys from a map
func mapKeys(rec map[string]interface{}) []string {
	keys := make([]string, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	return keys
}

func recordType(rec Record) string {
	keys := mapKeys(rec)
	if inList("arecord", keys) {
		return "arecord"
	} else if inList("notations", keys) {
		return "notation"
	} else if inList("verification_token", keys) {
		return "verification"
	} else if inList("input_values", keys) {
		return "input_values"
	} else if inList("presentation", keys) {
		return "presentation"
	} else if inList("urn", keys) {
		return "das"
	}
	return ""
}

func validateDasmaps(input string, verbose int) {
	if verbose > 0 {
		fmt.Println("validate dasmaps", input)
	}
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jsonBytes := scanner.Bytes()
		var rec Record
		err := json.Unmarshal(jsonBytes, &rec)
		if err != nil {
			log.Fatal(err)
		}
		switch t := recordType(rec); t {
		case "das":
			checkHash(rec)
		case "arecord":
			checkHash(rec)
		case "notation":
			checkHash(rec)
		case "verification":
			checkHash(rec)
		case "presentation":
			checkHash(rec)
		case "input_values":
			checkHash(rec)
		default:
			r, e := json.Marshal(rec)
			if e != nil {
				log.Fatal(e)
			}
			log.Fatal(fmt.Sprintf("unknown record type: %s %s", t, string(r)))
		}
	}

}

// validate keylearning das record
// {"keys": ["dataset", "file", "lumi"], "urn": "file_lumi4dataset", "system": "dbs3", "members": ["file.name", "lumi.number"]}
func validateKeylearning(input string, verbose int) {
	if verbose > 0 {
		fmt.Println("validate keylearning", input)
	}
	os.Exit(0)
}

func main() {
	var dasmaps string
	flag.StringVar(&dasmaps, "dasmaps", "", "dasmaps file")
	var keylearning string
	flag.StringVar(&keylearning, "keylearning", "", "keylearning file")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	if dasmaps != "" {
		validateDasmaps(dasmaps, verbose)
	} else if keylearning != "" {
		validateKeylearning(keylearning, verbose)
	}
}
