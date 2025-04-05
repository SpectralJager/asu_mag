// The control does not make full use of the ev3dev API where it could.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	TargetSpeed = flag.Float64("target", 100, "set target speed")
	P           = flag.Float64("p", .0, "set pid's p coefficient")
	I           = flag.Float64("i", .0, "set pid's i coefficient")
	D           = flag.Float64("d", .0, "set pid's d coefficient")
	Duration    = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"
	Sensor1    = "/sys/class/lego-sensor/sensor0"
	Sensor2    = "/sys/class/lego-sensor/sensor2"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	defer Command(LeftMotor, "stop")
	defer Command(RightMotor, "stop")

	pid1 := NewPIDController(*P, *I, *D, *TargetSpeed)
	pid2 := NewPIDController(*P, *I, *D, *TargetSpeed)

	after := time.After(time.Second * time.Duration(*Duration))

	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%f, Duration=%d sec\n", pid1.P, pid1.I, pid1.D, *TargetSpeed, *Duration)
	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%f, Duration=%d sec\n", pid2.P, pid2.I, pid2.D, *TargetSpeed, *Duration)

	SetSpeed(LeftMotor, int(0))
	SetSpeed(RightMotor, int(0))

	for {
		select {
		case <-after:
			return nil
		default:
		}
		leftSpeed, err := GetSpeed(LeftMotor)
		if err != nil {
			return fmt.Errorf("can't get speed of left motor: %w", err)
		}

		rightSpeed, err := GetSpeed(RightMotor)
		if err != nil {
			return fmt.Errorf("can't get speed of right motor: %w", err)
		}

		leftColor, err := GetIsBlack(Sensor1)
		if err != nil {
			return fmt.Errorf("can't get left color: %w", err)
		}

		rightColor, err := GetIsBlack(Sensor2)
		if err != nil {
			return fmt.Errorf("can't get right color: %w", err)
		}

		newSpeed_left := pid1.Update(float64(leftSpeed))
		newSpeed_right := pid2.Update(float64(rightSpeed))

		if leftColor <= 20 {
			newSpeed_right = *TargetSpeed
			newSpeed_left = 0
		} else if rightColor <= 20 {
			newSpeed_left = *TargetSpeed
			newSpeed_right = 0
		} else {
			newSpeed_left = *TargetSpeed
			newSpeed_right = *TargetSpeed
		}

		if newSpeed_left >= 1050 {
			newSpeed_left = 1050
		}
		if newSpeed_right >= 1050 {
			newSpeed_right = 1050
		}
		log.Printf("New speed: %0.2f, %0.2f\n", newSpeed_left, newSpeed_right)

		SetSpeed(LeftMotor, int(newSpeed_left))
		SetSpeed(RightMotor, int(newSpeed_right))
		Command(LeftMotor, "run-forever")
		Command(RightMotor, "run-forever")

		time.Sleep(time.Millisecond * 100)
	}
}

func Command(motor string, command string) error {
	commandPath := path.Join(motor, "command")

	return os.WriteFile(commandPath, []byte(command), 0644)
}

func GetIsBlack(sensor string) (int, error) {
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
	data, err := os.ReadFile(speedPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}
