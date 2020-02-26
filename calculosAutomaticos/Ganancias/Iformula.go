package Ganancias

type iformula interface {
	getResult() float64
	getResultInternal() float64
	getTope() *float64
	getNombre() string
}
