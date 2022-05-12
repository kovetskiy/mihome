package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/reconquest/executil-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Device struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Mac     string `json:"mac"`
	LocalIP string `json:"local_ip"`
	Token   string `json:"token"`
	Model   string `json:"model"`
}

func (device *Device) IsBulb() bool {
	return strings.HasPrefix(device.Model, "yeelink.light.color")
}

func (device *Device) IsVacuum() bool {
	return strings.HasPrefix(device.Model, "viomi.vacuum.")
}

func (device *Device) IsSwitch() bool {
	return strings.HasPrefix(device.Model, "chuangmi.plug.")
}

func (device *Device) IsLamp() bool {
	return strings.HasPrefix(device.Model, "yeelink.light.bslamp")
}

func getDevices(config *Config) ([]Device, error) {
	cmd := exec.Command("python3", config.Mi.Extractor)
	cmd.Env = []string{
		"MI_USERNAME=" + config.Mi.Username,
		"MI_PASSWORD=" + config.Mi.Password,
		"MI_SERVER=" + config.Mi.Server,
	}

	stdout, _, err := executil.Run(cmd)
	if err != nil {
		return nil, err
	}

	var devices []Device
	err = json.Unmarshal(stdout, &devices)
	if err != nil {
		return nil, karma.Format(err, "unmarshal json of stdout")
	}

	return devices, nil
}

func (device *Device) kind() string {
	switch {
	case device.IsBulb():
		return "yeelight"
	}

	panic(fmt.Sprintf("Unknown kind of device: %#v", device))
}

func (device *Device) command(command string, args ...string) error {
	params := []string{
		device.kind(),
		"--token", device.Token,
		"--ip", device.LocalIP,
		command,
	}

	params = append(params, args...)

	log.Debugf(nil, "%q", params)

	_, stderr, err := executil.Run(exec.Command("miiocli", params...))
	if err != nil {
		return err
	}

	if strings.Contains(string(stderr), "Error:") {
		return fmt.Errorf("command failed: %s", string(stderr))
	}

	return nil
}

func (device *Device) On() error {
	return device.command("on")
}

func (device *Device) Off() error {
	return device.command("off")
}

func (device *Device) SetRGB(rgb RGB) error {
	return device.command(
		"set_rgb",
		fmt.Sprint(rgb.Red),
		fmt.Sprint(rgb.Green),
		fmt.Sprint(rgb.Blue),
	)
}

func (device *Device) SetHex(hex string) error {
	rgb, err := HexToRGB(hex)
	if err != nil {
		return karma.Format(err, "convert hex to rgb")
	}

	return device.SetRGB(rgb)
}

func (device *Device) Toggle() error {
	return device.command("toggle")
}

func (device *Device) SetBrightness(level int) error {
	return device.command("set_brightness", fmt.Sprint(level))
}

func (device *Device) SetColorTemperature(level int) error {
	return device.command("set_color_temp", fmt.Sprint(level))
}
