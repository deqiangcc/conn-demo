package main

import (
	"bytes"
	"conn-demo/utils"
	"encoding/gob"
	"log"
	"time"
)



func main() {
	var gpssDecode []utils.GPS
	gpss := utils.NewGpss(100000)
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	encodeStart := time.Now()
	if err := enc.Encode(&gpss); err != nil {
		log.Fatal("encode err:", err)
	}
	log.Println("encode time:", time.Since(encodeStart))
	log.Println("data len:", len(buff.Bytes()))

	decodeStart := time.Now()
	dec := gob.NewDecoder(&buff)
	if err := dec.Decode(&gpssDecode); err != nil {
		log.Fatal("decode err:", err)
	}
	log.Println("decode time:", time.Since(decodeStart))
}