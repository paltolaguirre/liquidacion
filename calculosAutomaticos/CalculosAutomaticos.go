package calculosAutomaticos

import (
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type calculoAutomatico struct {
	Conceptoid  *int
	Liquidacion structLiquidacion.Liquidacion
}

func Calculoremunerativos() float64 {
	var importeCalculado float64

	return importeCalculado
}
