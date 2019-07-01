package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Row struct {
	K   string      `json:"k"`
	V   interface{} `json:"v"`
	T   string      `json:"t"`
	E   string      `json:"e"`
	TTL int         `json:"ttl"`
}

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
	_, err := p.Exec()

	if err != nil {
		log.Fatalln(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func BuildRedisOptions() *redis.Options {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		db = 0
	}
	opts := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
	return &opts
}

func BuildRedisClient(opts *redis.Options) *redis.Client {
	client := redis.NewClient(opts)
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	return client
}

func main() {
	log.Println("Started")

	rclient := BuildRedisClient(BuildRedisOptions())
	defer rclient.Close()

	rpipeline := rclient.Pipeline()
	defer rpipeline.Close()

	// Take file as input
	fn := os.Args[1]
	log.Println(fn)
	file, ferr := os.Open(fn)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := Row{}
		json.Unmarshal(scanner.Bytes(), &r)
		r.Load(rpipeline)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Completed")
}
