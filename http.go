package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/reconquest/pkg/log"
)

type Server struct {
	config  *Config
	devices []Device
}

func (server *Server) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	log.Debugf(nil, "%s %s", request.URL.String(), request.RemoteAddr)

	path := request.URL.Path
	if !strings.HasPrefix(path, "/mi/") {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	path = strings.TrimPrefix(path, "/mi")
	path = strings.TrimSuffix(path, "/")

	paramInt := func(key string) (int, bool) {
		value, err := strconv.Atoi(request.URL.Query().Get(key))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return 0, false
		}

		return value, true
	}

	paramString := func(key string) (string, bool) {
		value := request.URL.Query().Get(key)
		if value == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return "", false
		}

		return value, true
	}

	switch {
	case path == "/bulbs/on":
		server.controlBulbsOn()
	case path == "/bulbs/off":
		server.controlBulbsOff()
	case path == "/bulbs/toggle":
		server.controlBulbsToggle()
	case path == "/bulbs/brightness":
		if level, ok := paramInt("level"); ok {
			server.controlBulbsBrightness(level)
		}
	case path == "/bulbs/temperature":
		if level, ok := paramInt("level"); ok {
			server.controlBulbsColorTemperature(level)
		}
	case path == "/bulbs/color":
		if hex, ok := paramString("hex"); ok {
			server.controlBulbsHex(hex)
		}
	case path == "/bulbs/ukraine":
		server.controlBulbsModeUkraine()

	case path == "/bulbs/dance":
		server.controlBulbsModeDance()
	}
}

func (server *Server) bulbs() []Device {
	result := []Device{}
	for _, device := range server.devices {
		if device.IsBulb() {
			result = append(result, device)
		}
	}
	return result
}

func (server *Server) controlBulbs(
	fn func(id int, device Device) error,
) {
	workers := &sync.WaitGroup{}

	bulbs := server.bulbs()
	for id, bulb := range bulbs {
		workers.Add(1)
		go func(id int, bulb Device) {
			defer workers.Done()

			err := fn(id, bulb)
			if err != nil {
				log.Errorf(err, "control bulb %s", bulb.Name)
			}
		}(id+1, bulb)
	}

	workers.Wait()
}

func (server *Server) controlBulbsOn() {
	server.controlBulbs(func(_ int, device Device) error {
		return device.On()
	})
}

func (server *Server) controlBulbsOff() {
	server.controlBulbs(func(_ int, device Device) error {
		return device.Off()
	})
}

func (server *Server) controlBulbsToggle() {
	server.controlBulbs(func(_ int, device Device) error {
		return device.Toggle()
	})
}

func (server *Server) controlBulbsHex(hex string) {
	server.controlBulbs(func(_ int, device Device) error {
		return device.SetHex(hex)
	})
}

func (server *Server) controlBulbsBrightness(level int) {
	server.controlBulbs(func(_ int, device Device) error {
		return device.SetBrightness(level)
	})
}

func (server *Server) controlBulbsColorTemperature(level int) {
	server.controlBulbs(func(_ int, device Device) error {
		return device.SetColorTemperature(level)
	})
}

func (server *Server) controlBulbsModeUkraine() {
	yellow := "FFD500"
	blue := "005BBB"

	server.controlBulbs(func(id int, device Device) error {
		if id == 1 {
			return device.SetHex(yellow)
		}
		if id == 2 {
			return device.SetHex(yellow)
		}
		if id == 3 {
			return device.SetHex(blue)
		}

		return nil
	})
}

func (server *Server) controlBulbsModeDance() {
	colors := []string{
		"ff80ed",
		"00ffff",
		"00ff00",
		"0000ff",
		"ff0000",
		"ff00ff",
		"f6546a",
		"ff1493",
		"ff7373",
		"660066",
		"66cdaa",
		"fff68f",
		"ffd700",
	}

	randHex := func() string {
		return colors[rand.Intn(len(colors))]
	}

	started := time.Now()
	for time.Since(started).Seconds() < 20 {
		hex := randHex()
		server.controlBulbs(func(id int, device Device) error {
			return device.SetHex(hex)
		})
	}
}
