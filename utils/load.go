package utils

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// Row reprensents a row of the ru data files
type Row struct {
	K   string      `json:"k"`
	V   interface{} `json:"v"`
	T   string      `json:"t"`
	E   string      `json:"e"`
	TTL int         `json:"ttl"`
}

// ParseRow ...
func ParseRow(row []byte) *Row {
	r := Row{}
	json.Unmarshal(row, &r)
	return &r
}

// Load process the Row and loads it via a redis Pipeline
func (r *Row) Load(p redis.Pipeliner) {
	// log.Println(r.T)
	switch r.T {
	case "hash":
		p.HMSet(r.K, r.V.(map[string]interface{}))
	case "set":
		for _, el := range r.V.([]interface{}) {
			p.SAdd(r.K, el.(string))
		}
	case "zset":
		// {"k": "geo:venues:Yokohama Stadium",
		// "v": [["Baseball", 4171216862175648.0], ["Softball", 4171216862175648.0]],
		// "t": "zset", "ttl": -1}
		for _, el := range r.V.([]interface{}) {
			els := el.([]interface{})
			z := redis.Z{
				Member: els[0],
				Score:  els[1].(float64),
			}
			p.ZAdd(r.K, z)
		}
	case "list":
		for _, el := range r.V.([]interface{}) {
			p.RPush(r.K, el.(string))
		}
	case "string":
		if r.E == "embstr" {
			p.Set(r.K, r.V.(string), -1)
		} else if r.E == "raw" {
			ary, _ := base64.StdEncoding.DecodeString(r.V.(string))
			for i, el := range ary {
				// BITFIELD is not supported yet by redis-go
				script := fmt.Sprintf(
					`return redis.call("BITFIELD", "%s", "SET", "u8", "%d", "%d")`,
					r.K, i*8, el,
				)
				p.Eval(script, []string{})
			}
		}
	default:
		log.Printf("Don't know how to process %s", r.T)
	}
	if r.TTL > 0 {
		p.Expire(r.K, time.Duration(r.TTL))
	}
}

// LoadFromFile ...
func LoadFromFile(fn string) {
	file, ferr := os.Open(fn)
	if ferr != nil {
		log.Fatal(ferr)
		return
	}
	defer file.Close()

	rclient := BuildRedisClient(BuildRedisOptions(), true)
	defer rclient.Close()

	rpipeline := rclient.Pipeline()
	defer rpipeline.Close()

	batchSize := 100
	counter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ParseRow(scanner.Bytes()).Load(rpipeline)
		counter++
		if counter%batchSize == 0 {
			_, err := rpipeline.Exec()

			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	_, err := rpipeline.Exec()

	if err != nil {
		log.Fatalln(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("total keys loaded: %d", counter)
}
