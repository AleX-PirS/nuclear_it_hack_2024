package handlers

import (
	"log"

	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/http/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/paulmach/orb/geojson"
)
	

type Handler struct {
	outCh chan *dto.Request
	respCh chan *geojson.FeatureCollection
}

func New() *Handler{
	return &Handler{
		outCh: make(chan *dto.Request),
		respCh: make(chan *geojson.FeatureCollection),
	}
}

func (h *Handler) GetChans() (chan *dto.Request, chan *geojson.FeatureCollection) {
	return h.outCh, h.respCh
}

func (h *Handler) HandleJsons(c *fiber.Ctx) error {
	jsonData := &dto.Request{}
	if err := c.BodyParser(jsonData); err != nil{
		return c.SendStatus(500)
	}
	
	h.outCh <- jsonData

	geo := <- h.respCh
	data, err := geo.MarshalJSON()
	if err != nil {
		log.Fatal("Error marshall json")
	}
	return c.JSON(data)
}