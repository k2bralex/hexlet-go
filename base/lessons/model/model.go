package model

import "math"

type (
	AreaCalc interface {
		Area() float64
	}

	Rectangle struct {
		A, B, C int64
	}

	Circle struct {
		R int
	}

	Calculator struct {
		AreaCalc
	}
)

func (r *Rectangle) Area() float64 {
	return 2.33
}

func (c *Circle) Area() float64 {
	return math.Pi * float64(c.R) * float64(c.R)
}

func CalcArea(s AreaCalc) float64 {
	return s.Area()
}
