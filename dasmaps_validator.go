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
	"runtime"
	"time"
)

// Record represents DAS map record
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

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var input string
	flag.StringVar(&input, "input", "", "input file")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	if version {
		fmt.Println(info())
		return
	}
	validateDasmaps(input, verbose)
}
