package models

import (
	"encoding/json"
	"fmt"
)

type Price int32

func NewPrice(f float64) Price {
	return Price((f * 100) + 0.5)
}

func (m Price) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

func (m Price) Multiply(f float64) Price {
	x := (float64(m) * f) + 0.5
	return Price(x)
}

func (m Price) String() string {
	x := float64(m)
	x = x / 100
	return fmt.Sprintf("$%.2f", x)
}

func (m Price) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Float64())
}
