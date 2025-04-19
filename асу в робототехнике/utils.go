package asu

import (
	"log"
	"os"
	"path"
	"strconv"
)

func Must[T any](val T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func Command(motor string, command string) error {
	commandPath := path.Join(motor, "command")

	return os.WriteFile(commandPath, []byte(command), 0644)
}

func GetDistance(sensor string) (int, error) {
	distancePath := path.Join(sensor, "value0")
	data, err := os.ReadFile(distancePath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}

func SetSpeed(motor string, speed int, lb, ub int) error {
	if speed >= ub {
		speed = ub
	}
	if speed <= lb {
		speed = -lb
	}
	speedPath := path.Join(motor, "speed_sp")

	return os.WriteFile(speedPath, []byte(strconv.Itoa(speed)), 0644)
}

func GetSpeed(motor string) (int, error) {
	speedPath := path.Join(motor, "speed")
	// tic/sec
	data, err := os.ReadFile(speedPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}

func GetColor(sensor string) (int, error) {
	distancePath := path.Join(sensor, "value0")
	data, err := os.ReadFile(distancePath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}
