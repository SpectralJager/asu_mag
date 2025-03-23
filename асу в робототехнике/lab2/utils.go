package main

import (
	"log"

	"github.com/ev3go/ev3dev"
)

func Must[T any](val T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func SetMotorPID(motor *ev3dev.TachoMotor, P, I, D int) error {
	motor.SetHoldPIDKp(P).
		SetHoldPIDKi(I).
		SetHoldPIDKd(D)
	return nil
}
