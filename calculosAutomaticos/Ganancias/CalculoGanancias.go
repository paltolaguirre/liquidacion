package Ganancias

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
)

type CalculoGanancias struct {
	Liquidacionitem *structLiquidacion.Liquidacionitem
	Liquidacion     *structLiquidacion.Liquidacion
	Db              *gorm.DB
	EjecutarCalculo bool
}

func (cg *CalculoGanancias) getResultOnDemandTemplate(codigo string, orden int, formula iformula) float64 {

	var importeTotal float64
	var importePuntero *float64
	var topePuntero *float64

	for _, acumulador := range cg.Liquidacionitem.Acumuladores {
		if acumulador.Codigo == codigo {
			importePuntero = acumulador.Importe
			topePuntero = acumulador.Tope
		}
	}

	if importePuntero == nil {
		if cg.EjecutarCalculo == false {
			panic(errors.New("No se pudo obtener el valor de" + formula.getNombre() + " para la liquidacion con mes de liquidacion " + cg.Liquidacion.Fechaperiodoliquidacion.Month().String()))
		}
		importeTotal = formula.getResultInternal()
		topePuntero = formula.getTope()
		fmt.Println("Calculos Automaticos -", formula.getNombre()+":", importeTotal)
		importePuntero = &importeTotal
		acumuladorRembruta := structLiquidacion.Acumulador{
			Nombre:      formula.getNombre(),
			Codigo:      codigo,
			Descripcion: "",
			Orden:       orden,
			Importe:     importePuntero,
			Tope:        topePuntero,
			Esmostrable: formula.getEsMostrable(),
		}
		cg.Liquidacionitem.Acumuladores = append(cg.Liquidacionitem.Acumuladores, acumuladorRembruta)
	} else {
		importeTotal = *importePuntero
	}
	return importeTotal
}

func (cg *CalculoGanancias) Calculate() float64 {
	cg.invocarCalculosLiquidacionAnual()
	return (&CalculoRetencionDelMes{*cg}).getResult()
}

func (cg *CalculoGanancias) invocarCalculosLiquidacionAnual() {
	(&CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras{*cg}).getResult()
	(&CalculoHorasExtrasRemuneracionExenta{*cg}).getResult()
	(&CalculoMovilidadYViaticosRemuneracionExenta{*cg}).getResult()
	(&CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta{*cg}).getResult()
	(&CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos{*cg}).getResult()
	(&CalculoHorasExtrasRemuneracionExentaOtrosEmpleos{*cg}).getResult()
	(&CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos{*cg}).getResult()
	(&CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos{*cg}).getResult()
	(&CalculoSubtotalRemuneracionGravada{*cg}).getResult()
	(&CalculoSubtotalRemuneracionNoGravadaNoAlcanzadaExenta{*cg}).getResult()
	(&CalculoTotalRemuneraciones{*cg}).getResult()
	(&CalculoPrimasDeSeguroParaElCasoDeMuerteAnual{*cg}).getResult()
	(&CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual{*cg}).getResult()
	(&CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual{*cg}).getResult()
	(&CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro{*cg}).getResult()
	(&CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica{*cg}).getResult()
	(&CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares{*cg}).getResult()
	(&CalculoSubtotalDeduccionesGenerales{*cg}).getResult()
	(&CalculoSubtotalAnual{*cg}).getResult()
	(&CalculoGananciaNetaAnual{*cg}).getResult()
	(&CalculoConyugeAnual{*cg}).getResult()
	(&CalculoHijosAnual{*cg}).getResult()
	(&CalculoSubtotalCargasFamilia{*cg}).getResult()
	(&CalculoSubtotalDeduccionesPersonalesAnual{*cg}).getResult()
	(&CalculoRemuneracionSujetaAImpuesto{*cg}).getResult()
	(&CalculoRemuneracionSujetaAImpuestoSinIncluirHorasExtras{*cg}).getResult()
	(&CalculoAlicuotaArt90LeyGanancias{*cg}).getResult()
	(&CalculoAlicuotaAplicableSinIncluirHorasExtras{*cg}).getResult()
	(&CalculoImpuestoDeterminado{*cg}).getResult()
	(&CalculoPagosACuenta{*cg}).getResult()
	(&CalculoSaldoAPagar{*cg}).getResult()
}

func (cg *CalculoGanancias) getfgSacCuotas(correspondeSemestre bool) float64 {
	var importeTotal, importeConcepto float64

	if correspondeSemestre {
		for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
			liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
			concepto := liquidacionitem.Concepto
			var mes float64 = 1
			if concepto.Basesac == true {
				if concepto.Prorrateo == true {
					mes = float64(cg.getfgMesesAProrratear(concepto))
				}
				importeLiquidacionitem := liquidacionitem.Importeunitario
				if importeLiquidacionitem != nil {
					importeConcepto = *importeLiquidacionitem / mes
				}

				if *concepto.Tipoconceptoid == -4 {
					importeConcepto = importeConcepto * -1
				}
				importeTotal = importeTotal + importeConcepto
			}
		}

		importeTotal = importeTotal + cg.getfgBaseSacOtrosEmpleos()

		importeTotal = importeTotal + cg.obtenerConceptosProrrateoMesesAnteriores()
	}

	return importeTotal / 12
}

type importeMes struct {
	Importeunitario *float64
	Mesliquidacion  string
}

func (cg *CalculoGanancias) obtenerConceptosProrrateoMesesAnteriores() float64 {
	var importemes []importeMes
	anioLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("01")
	legajoID := cg.Liquidacion.Legajoid

	sql := "SELECT li.importeunitario, to_char(l.fechaperiodoliquidacion, 'MM') AS mesliquidacion FROM liquidacion l INNER JOIN liquidacionitem li on l.id = li.liquidacionid INNER JOIN legajo le on le.id = l.legajoid INNER JOIN concepto c on c.id = li.conceptoid WHERE li.ID != " + strconv.Itoa(cg.Liquidacion.ID) + " AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioLiquidacion) + "' AND to_char(l.fechaperiodoliquidacion, 'MM') < '" + mesLiquidacion + "' AND le.id = " + strconv.Itoa(*legajoID) + " and c.prorrateo = true ORDER BY to_char(l.fechaperiodoliquidacion, 'MM') ASC"
	cg.Db.Raw(sql).Scan(&importemes)
	var trece float64 = 13
	var importeTotal float64 = 0
	for i := 0; i < len(importemes); i++ {
		mesLiquidacion, _ := strconv.ParseFloat(importemes[i].Mesliquidacion, 64)
		importeConcepto := *importemes[i].Importeunitario / (trece - mesLiquidacion)

		importeTotal = importeTotal + importeConcepto
	}

	return importeTotal

}

func (cg *CalculoGanancias) getfgBaseSacOtrosEmpleos() float64 {

	var arrayBaseSacPositivos []float64
	var arrayBaseSacNegativos []float64

	arrayBaseSacPositivos = append(arrayBaseSacPositivos, (&CalculoRemuneracionBrutaOtrosEmpleos{*cg}).getResult())
	arrayBaseSacPositivos = append(arrayBaseSacPositivos, (&CalculoRemuneracionNoHabitualOtrosEmpleos{*cg}).getResult())
	arrayBaseSacPositivos = append(arrayBaseSacPositivos, (&CalculoHorasExtrasGravadasOtrosEmpleos{*cg}).getResult())
	arrayBaseSacPositivos = append(arrayBaseSacPositivos, (&CalculoMovilidadYViaticosGravadaOtrosEmpleos{*cg}).getResult())
	arrayBaseSacPositivos = append(arrayBaseSacPositivos, (&CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos{*cg}).getResult())

	arrayBaseSacNegativos = append(arrayBaseSacNegativos, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos{*cg}).getResult())
	arrayBaseSacNegativos = append(arrayBaseSacNegativos, (&CalculoAportesObraSocialOtrosEmpleos{*cg}).getResult())
	arrayBaseSacNegativos = append(arrayBaseSacNegativos, (&CalculoCuotaSindicalOtrosEmpleos{*cg}).getResult())

	return Sum(arrayBaseSacPositivos) - Sum(arrayBaseSacNegativos)
}

func (cg *CalculoGanancias) obtenerRemunerativosOtros() float64 {
	var arrayRemunerativosOtros []float64
	var totalRemunerativosOtros float64

	arrayRemunerativosOtros = append(arrayRemunerativosOtros, (&CalculoRemuneracionBrutaOtrosEmpleos{*cg}).getResult())
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, (&CalculoSACPrimerCuotaOtrosEmpleos{*cg}).getResult())
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, (&CalculoSACSegundaCuotaOtrosEmpleos{*cg}).getResult())
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, (&CalculoHorasExtrasGravadasOtrosEmpleos{*cg}).getResult())

	totalRemunerativosOtros = Sum(arrayRemunerativosOtros)
	fmt.Println("Calculos Automaticos - Remunerativos Otros:", totalRemunerativosOtros)
	return totalRemunerativosOtros
}

func (cg *CalculoGanancias) obtenerRemunerativosMenosDescuentos() float64 {
	var totalRemunerativos, totalDescuentos float64
	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
		tipoconcepto := *liquidacionitem.Concepto.Tipoconceptoid
		importeconcepto := liquidacionitem.Importeunitario
		if importeconcepto != nil {

			if tipoconcepto == -1 {
				totalRemunerativos = totalRemunerativos + *importeconcepto
			}
			if tipoconcepto == -3 {
				totalDescuentos = totalDescuentos + *importeconcepto
			}
		}
	}
	fmt.Println("Calculos Automaticos - RemunerativosMenosDescuentos:", totalRemunerativos-totalDescuentos)
	return totalRemunerativos - totalDescuentos
}

func (cg *CalculoGanancias) GetfgImporteTotalSegunTipoImpuestoGanancias(tipoImpuestoALasGanancias string) float64 {
	var importeTotal, importeConcepto float64

	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		tipoimpuesto := obtenerTipoImpuesto(concepto, cg.Db)
		var mes float64 = 1

		if tipoimpuesto == tipoImpuestoALasGanancias && concepto.Codigo != "IMPUESTO_GANANCIAS" && concepto.Codigo != "IMPUESTO_GANANCIAS_DEVOLUCION" {
			if concepto.Prorrateo == true {
				mes = float64(cg.getfgMesesAProrratear(concepto))
			}
			importeLiquidacionitem := liquidacionitem.Importeunitario
			if importeLiquidacionitem != nil {
				if concepto.ID == -6 {
					importeConcepto = (*importeLiquidacionitem / float64(2)) / mes
				} else {
					importeConcepto = *importeLiquidacionitem / mes
				}

			}
			importeTotal = importeTotal + importeConcepto

		}
	}
	return importeTotal
}

func (cg *CalculoGanancias) getfgMesesAProrratear(concepto *structConcepto.Concepto) int {
	fechadesde := strconv.Itoa(cg.Liquidacion.Fechaperiodoliquidacion.Year()) + "-01-01"
	fechahasta := cg.Liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")
	var fechaliquidacionmasantigua *time.Time
	var sql string
	mesAProrratear := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)

	sql = "SELECT l.fechaperiodoliquidacion FROM Liquidacion l INNER JOIN Liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid INNER JOIN legajo le ON le.id = l.legajoid WHERE c.id = " + strconv.Itoa(concepto.ID) + " AND l.fechaperiodoliquidacion BETWEEN '" + fechadesde + "' AND '" + fechahasta + "' AND le.id = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + "AND le.deleted_at IS NULL AND l.deleted_at IS NULL and li.deleted_at IS NULL, c.deleted_at IS NULL ORDER BY fechaperiodoliquidacion ASC LIMIT 1"
	cg.Db.Raw(sql).Row().Scan(&fechaliquidacionmasantigua)

	if fechaliquidacionmasantigua != nil {
		mesLiquidacionBD := getfgMes(fechaliquidacionmasantigua)

		if mesLiquidacionBD < mesAProrratear {
			mesAProrratear = mesLiquidacionBD
		}
	}
	fmt.Println("Calculos Automaticos - Mes a Prorratear:", 13-mesAProrratear)
	return 13 - mesAProrratear

}

func (cg *CalculoGanancias) getfgDetalleCargoFamiliar(columnaDetalleCargoFamiliar string, valorfijocolumna string, porcentaje float64) float64 {
	var importeTotal float64
	var tienevalorbeneficio bool
	anioperiodoliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND s.deleted_at IS NULL AND dcfs.deleted_at IS NULL"
	cg.Db.Raw(sql).Scan(&detallecargofamiliar)
	sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE to_number(to_char(bs.mesdesde, 'MM'),'99') <= " + strconv.Itoa(mesperiodoliquidacion) + " AND to_number(to_char(bs.meshasta, 'MM'), '99') > " + strconv.Itoa(mesperiodoliquidacion) + " AND bs.siradigtipogrillaid = -24 AND to_char(s.periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND s.deleted_at IS NULL AND bs.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&tienevalorbeneficio)

	if detallecargofamiliar.ID != 0 {

		mesdadobaja := getfgMes(detallecargofamiliar.Meshasta)
		mesdadoalta := getfgMes(detallecargofamiliar.Mesdesde)
		valorfijo := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", valorfijocolumna)

		if tienevalorbeneficio == true {
			valorfijo = valorfijo * 1.22
		}

		if mesdadobaja == 0 {
			importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
		} else {
			if mesdadobaja <= mesperiodoliquidacion {
				importeTotal = (valorfijo / 12) * float64(mesdadobaja-mesdadoalta) * (porcentaje / 100)
			} else {
				if mesdadobaja > mesperiodoliquidacion {
					importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
				} else {
					if mesdadoalta > mesperiodoliquidacion {
						importeTotal = 0
					}
				}
			}
		}
	}

	return importeTotal
}

func (cg *CalculoGanancias) getfgDetalleCargoFamiliarAnual(columnaDetalleCargoFamiliar string, valorfijocolumna string, porcentaje float64, valorfijoMNI float64) float64 {
	var importeTotal float64
	var tienevalorbeneficio bool
	anioperiodoliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND s.deleted_at IS NULL AND dcfs.deleted_at IS NULL"
	cg.Db.Raw(sql).Scan(&detallecargofamiliar)
	sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE to_number(to_char(bs.mesdesde, 'MM'),'99') <= " + strconv.Itoa(mesperiodoliquidacion) + " AND to_number(to_char(bs.meshasta, 'MM'), '99') > " + strconv.Itoa(mesperiodoliquidacion) + " AND bs.siradigtipogrillaid = -24 AND to_char(s.periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND s.deleted_at IS NULL AND bs.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&tienevalorbeneficio)

	if detallecargofamiliar.ID != 0 {

		if *detallecargofamiliar.Montoanual < valorfijoMNI {

			mesdadobaja := getfgMes(detallecargofamiliar.Meshasta)
			mesdadoalta := getfgMes(detallecargofamiliar.Mesdesde)
			valorfijo := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", valorfijocolumna)

			if tienevalorbeneficio == true {
				valorfijo = valorfijo * 1.22
			}

			if mesdadobaja == 0 {
				importeTotal = (valorfijo / 12) * -float64(12-mesdadoalta) * (porcentaje / 100)
			} else {
				if mesdadobaja <= 12 {
					importeTotal = (valorfijo / 12) * float64(mesdadobaja-mesdadoalta) * (porcentaje / 100)
				}
			}
		}
	}

	return importeTotal
}

func (cg *CalculoGanancias) getfgValorFijoImpuestoGanancia(nombretabla string, nombrecolumna string) float64 {
	var importeTope float64
	anioLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	sql := "SELECT " + nombrecolumna + " FROM " + nombretabla + " WHERE anio = " + strconv.Itoa(anioLiquidacion)
	cg.Db.Raw(sql).Row().Scan(&importeTope)

	return importeTope
}

func (cg *CalculoGanancias) obtenerLiquidacionesIgualAnioLegajoMenorMes() *[]structLiquidacion.Liquidacion {
	var liquidaciones []structLiquidacion.Liquidacion
	anioperiodoliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
	cg.Db.Order("to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') desc").Set("gorm:auto_preload", true).Find(&liquidaciones, "to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') < ? AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND legajoid = ?", mesliquidacion, strconv.Itoa(anioperiodoliquidacion), *cg.Liquidacion.Legajoid)

	return &liquidaciones
}

func (cg *CalculoGanancias) obtenerLiquidacionIgualAnioLegajoMesAnterior() *structLiquidacion.Liquidacion {
	var liquidacionMesAnterior structLiquidacion.Liquidacion
	liquidaciones := *cg.obtenerLiquidacionesIgualAnioLegajoMenorMes()
	if len(liquidaciones) > 0 {
		liquidacionMesAnterior = liquidaciones[0]
	}
	return &liquidacionMesAnterior
}

func (cg *CalculoGanancias) getfgImporteGananciasOtroEmpleoSiradig(columnaimportegananciasotroempleosiradig string) float64 {
	var importeTotal float64
	anoLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("2006")
	mesLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("01")
	legajoid := strconv.Itoa(*cg.Liquidacion.Legajoid)
	sql := "SELECT SUM(" + columnaimportegananciasotroempleosiradig + ") FROM importegananciasotroempleosiradig WHERE '" + anoLiquidacion + "' = extract(YEAR from mes) and '" + mesLiquidacion + "' = extract(MONTH from mes) " +
		"and siradigid in (SELECT id from siradig where legajoid = " + legajoid + " ) AND importegananciasotroempleosiradig.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func (cg *CalculoGanancias) getfgImporteTotalSiradigSegunTipoGrillaSinMes(columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string) float64 {
	var importeTotal float64
	anioliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = ts.siradigid WHERE stg.codigo = '" + tipodeducciondesgravacionsiradig + "' AND sdg.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anioliquidacion) + "' AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL AND ts.deleted_at IS NULL;"
	cg.Db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func (cg *CalculoGanancias) getfgImporteTotalSiradigSegunTipoGrilla(columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string) float64 {
	var importeTotal float64
	mesLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("01")
	anioliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = ts.siradigid WHERE to_number(to_char(mes, 'MM'),'99') <= " + mesLiquidacion + " AND stg.codigo = '" + tipodeducciondesgravacionsiradig + "' AND sdg.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anioliquidacion) + "' AND ts.deleted_at IS  NULL AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL;"
	cg.Db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}
