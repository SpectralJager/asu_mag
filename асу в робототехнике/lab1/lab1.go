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
	TargetSpeed = flag.Int("target", 200, "set target speed")
	P           = flag.Float64("p", .0, "set pid's p coefficient")
	I           = flag.Float64("i", .0, "set pid's i coefficient")
	D           = flag.Float64("d", .0, "set pid's d coefficient")
	Duration    = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"
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

	pid1 := NewPIDController(*P, *I, *D, float64(*TargetSpeed))
	pid2 := NewPIDController(*P, *I, *D, float64(*TargetSpeed))

	after := time.After(time.Second * time.Duration(*Duration))

	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%d, Duration=%d sec\n", pid1.P, pid1.I, pid1.D, *TargetSpeed, *Duration)
	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%d, Duration=%d sec\n", pid2.P, pid2.I, pid2.D, *TargetSpeed, *Duration)

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

		newSpeed_left := pid1.Update(float64(leftSpeed))
		newSpeed_right := pid2.Update(float64(rightSpeed))
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
