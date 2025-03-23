// The control does not make full use of the ev3dev API where it could.
package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	Target   = flag.Int("target", 100, "set target distance")
	P        = flag.Float64("p", .0, "set pid's p coefficient")
	I        = flag.Float64("i", .0, "set pid's i coefficient")
	D        = flag.Float64("d", .0, "set pid's d coefficient")
	Step     = flag.Float64("step", 50., "set speed step increment")
	Duration = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"
	Sensor     = "/sys/class/lego-sensor/sensor0"
)

func main() {
	flag.Parse()
	if *Target < 100 {
		*Target = 200
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	defer Command(LeftMotor, "stop")
	defer Command(RightMotor, "stop")

	pid := NewPIDController(*P, *I, *D, float64(*Target))

	SetSpeed(LeftMotor, 1)
	SetSpeed(RightMotor, 1)

	after := time.After(time.Second * time.Duration(*Duration))
	for {
		select {
		case <-after:
			return nil
		default:
		}
		distance, err := GetDistance(Sensor)
		if err != nil {
			return err
		}
		newDistance := pid.Update(float64(distance))
		k := newDistance / pid.Target * -1

		speed, err := GetSpeed(LeftMotor)
		if err != nil {
			return err
		}
		newSpeed := speed
		switch {
		// case distance < int(pid.Target/3) || distance < 100:
		// 	newSpeed = 0
		case k > 0 || k < 0:
			newSpeed += int(*Step * k)
		}

		SetSpeed(LeftMotor, newSpeed)
		SetSpeed(RightMotor, newSpeed)

		Command(LeftMotor, "run-forever")
		Command(RightMotor, "run-forever")

		log.Printf("distance: %d mm\nspeed: %d\nk = %.2f\n", distance, newSpeed, k)

		time.Sleep(time.Millisecond * 100)
	}
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

func SetSpeed(motor string, speed int) error {
	if speed >= 1050 {
		speed = 1050
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
