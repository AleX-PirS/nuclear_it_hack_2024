package handlers

import (
	"encoding/json"

	"github.com/AleX-PirS/nuclear_it_hack_2024/interfaces/http/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)
	

type Handler struct {
	outCh chan dto.Request
	respCh chan dto.Response
}

func New() *Handler{
	return &Handler{
		outCh: make(chan dto.Request),
		respCh: make(chan dto.Response),
	}
}

func (h *Handler) GetChans() (chan dto.Request, chan dto.Response) {
	return h.outCh, h.respCh
}


func (h *Handler) HandleJsons(c *fiber.Ctx) error {
	log.Info("New json.")
	jsonData := dto.Request{}
	var inpData map[string]interface{}
	err := c.BodyParser(&inpData)
	if err != nil {
		log.Warn("BAD!", err.Error())
	}
	log.Info(inpData)
	err = c.BodyParser(&jsonData)
	if err != nil {
		log.Warn("BAD!", err.Error())
	}

	log.Info(jsonData)

	// if err := c.BodyParser(&jsonData); err != nil{
	// 	log.Info(jsonData)
	// 	return c.SendStatus(500)
	// }
	h.outCh <- dto.Request{Accuracy: 1}
	// h.outCh <- jsonData

	geo := <- h.respCh
	data, err := json.Marshal(geo)
	if err != nil {
		log.Fatal("Error marshall json")
	}
	return c.JSON(data)
}