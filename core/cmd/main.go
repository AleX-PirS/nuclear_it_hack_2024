package main

import (
	"log"

	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/handlers"
	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/http"
	"github.com/paulmach/orb/geojson"
)

func main(){
	fApp := http.New()
	h := handlers.New()
	s := http.NewServer(fApp, h)
	readCh, sendCh := s.GetChans()

	go s.ConfigurateAndRun()

	for {
		data := <-readCh
		log.Println(data)
		fc := geojson.NewFeatureCollection()
		fc.Type = "TESTING"
		sendCh <- fc
	}
}
