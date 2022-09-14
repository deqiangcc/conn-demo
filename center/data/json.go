package main

import (
	"conn-demo/utils"
	"encoding/json"
	"log"
	"time"
)

func main() {
	var gpssDecode []utils.GPS
	gpss := utils.NewGpss(100000)

	startEncode := time.Now()
	data, err := json.Marshal(gpss)
	if err != nil {
		log.Fatal("encode err:", err)
	}
	log.Println("encode time:", time.Since(startEncode))
	log.Println("encode data len:", len(data))

	startDecode := time.Now()
	if err := json.Unmarshal(data, &gpssDecode); err != nil {
		log.Fatal("decode err:", err)
	}
	log.Println("decode time:", time.Since(startDecode))
}
