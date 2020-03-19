package Ganancias

type CalculoSaldoAPagar struct {
	CalculoGanancias
}

func (cg *CalculoSaldoAPagar) getResultInternal() float64 {

	return ((&CalculoImpuestoDeterminado{cg.CalculoGanancias}).getResult() - (&CalculoRetencionAcumulada{cg.CalculoGanancias}).getResult() - (&CalculoPagosACuenta{cg.CalculoGanancias}).getResult())

}

func (cg *CalculoSaldoAPagar) getResult() float64 {
	return cg.getResultOnDemandTemplate("SALDO_A_PAGAR", 0, cg)
}

func (cg *CalculoSaldoAPagar) getTope() *float64 {
	return nil
}

func (cg *CalculoSaldoAPagar) getNombre() string {
	return "Saldo a Pagar"
}

func (cg *CalculoSaldoAPagar) getEsMostrable() bool {
	return false
}
