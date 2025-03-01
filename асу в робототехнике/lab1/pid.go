package main

type PIDController struct {
	P      float64
	I      float64
	D      float64
	Target float64
	State  PIDState
}

type PIDState struct {
	Error           float64
	ErrorIntegral   float64
	ErrorDerivative float64
}

func NewPIDController(P, I, D, target float64) PIDController {
	return PIDController{
		P:      P,
		I:      I,
		D:      D,
		Target: target,
	}
}

func (pid *PIDController) Update(input float64) float64 {
	prevError := pid.State.Error
	pid.State.Error = pid.Target - input
	pid.State.ErrorIntegral += pid.State.Error
	pid.State.ErrorDerivative = prevError - pid.State.Error

	return pid.P*pid.State.Error +
		pid.I*pid.State.ErrorIntegral +
		pid.D*pid.State.ErrorDerivative
}
