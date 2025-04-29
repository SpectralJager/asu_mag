// The control does not make full use of the ev3dev API where it could.
package main

import (
	"asu"
	"flag"
	"log"
	"time"
)

var (
	P        = flag.Float64("p", .0, "set pid's p coefficient")
	I        = flag.Float64("i", .0, "set pid's i coefficient")
	D        = flag.Float64("d", .0, "set pid's d coefficient")
	Duration = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"

	DistanceSensor = "/sys/class/lego-sensor/sensor4"
	ColorSensor1   = "/sys/class/lego-sensor/sensor2"
	ColorSensor2   = "/sys/class/lego-sensor/sensor3"
	GyroSensor     = "/sys/class/lego-sensor/sensor5"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	defer asu.Command(LeftMotor, "stop")
	defer asu.Command(RightMotor, "stop")

	after := time.After(time.Second * time.Duration(*Duration))

	// correctionPID := asu.NewPIDController(.0, .0, .0, 0)

	asu.SetSpeed(LeftMotor, int(0), 0, 0)
	asu.SetSpeed(RightMotor, int(0), 0, 0)

	state := 0
	initAngle := 0
	counter := 0
	for {
		select {
		case <-after:
			return nil
		default:
		}

		log.Printf("state: %d, ", state)
		switch state {
		case 0:
			// follow line
			distance, _ := asu.GetDistance(DistanceSensor)
			log.Printf("distance: %d\n", distance)
			if distance < 200 {
				state = 1
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 1:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 80 || dif < -80 {
				state = 2
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, -100, -160, 160)
				asu.SetSpeed(RightMotor, 100, -160, 160)
			}
		case 2:
			counter += 1
			if counter == 50 {
				counter = 0
				state = 3
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 3:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 80 || dif < -80 {
				state = 4
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, -160, 160)
				asu.SetSpeed(RightMotor, -100, -160, 160)
			}
		case 4:
			counter += 1
			if counter == 100 {
				counter = 0
				state = 5
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 5:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 80 || dif < -80 {
				state = 6
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, -160, 160)
				asu.SetSpeed(RightMotor, -100, -160, 160)
			}
		case 6:
			counter += 1
			if counter == 50 {
				counter = 0
				state = 7
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 7:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 80 || dif < -80 {
				state = 1
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, -100, -160, 160)
				asu.SetSpeed(RightMotor, 100, -160, 160)
			}
		}

		// leftColor, err := asu.GetColor(ColorSensor1)
		// if err != nil {
		// 	return fmt.Errorf("can't get left color: %w", err)
		// }
		// rightColor, err := asu.GetColor(ColorSensor2)
		// if err != nil {
		// 	return fmt.Errorf("can't get right color: %w", err)
		// }
		// correction := correctionPID.Update(float64(leftColor - rightColor))

		// newLeftSpeed := int(80 + (-1 * correction))
		// newRightSpeed := int(80 + correction)

		// distance, err := asu.GetDistance(DistanceSensor)
		// if err != nil {
		// 	return err
		// }

		// if distance < 200 {

		// } else {
		// }

		asu.Command(LeftMotor, "run-forever")
		asu.Command(RightMotor, "run-forever")
		time.Sleep(time.Millisecond * 100)
	}
}
