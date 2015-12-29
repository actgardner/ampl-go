package model

type Shape int

const (
	Constant Shape = iota
	Linear
	Quadratic
	NonLinear
)

func (c Shape) String() string {
	switch c {
	case Constant:
		return "Constant"
	case Linear:
		return "Linear"
	case Quadratic:
		return "Quadratic"
	case NonLinear:
		return "Non-Linear"
	}
	return "Unknown"
}
