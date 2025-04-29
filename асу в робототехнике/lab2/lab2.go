// The control does not make full use of the ev3dev API where it could.
package main

import (
	"asu"
	"flag"
	"log"
	"time"
)

var (
	Target   = flag.Int("target", 100, "set target distance")
	P        = flag.Float64("p", .0, "set pid's p coefficient")
	I        = flag.Float64("i", .0, "set pid's i coefficient")
	D        = flag.Float64("d", .0, "set pid's d coefficient")
	Step     = flag.Float64("step", 10., "set speed step increment")
	Duration = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"
	Sensor     = "/sys/class/lego-sensor/sensor1"
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
	defer asu.Command(LeftMotor, "stop")
	defer asu.Command(RightMotor, "stop")

	pid := asu.NewPIDController(*P, *I, *D, float64(*Target))

	asu.SetSpeed(LeftMotor, 0, 0, 0)
	asu.SetSpeed(RightMotor, 0, 0, 0)

	after := time.After(time.Second * time.Duration(*Duration))
	distances := []int{}
	for {
		select {
		case <-after:
			return nil
		default:
		}
		distance, err := asu.GetDistance(Sensor)
		if err != nil {
			return err
		}
		distances = append(distances, distance)
		if len(distances) > 3 {
			total := 0
			for _, d := range distances[len(distances)-3:] {
				total += d
			}
			distance = total / 4
		}
		newDistance := int(pid.Update(float64(distance)))

		newSpeed := int(-newDistance)

		asu.SetSpeed(LeftMotor, newSpeed, -1050, 1050)
		asu.SetSpeed(RightMotor, newSpeed, -1050, 1050)

		asu.Command(LeftMotor, "run-forever")
		asu.Command(RightMotor, "run-forever")

		log.Printf("distance: %d mm, speed: %d\n", distance, newSpeed)

		time.Sleep(time.Millisecond * 200)
	}
}
