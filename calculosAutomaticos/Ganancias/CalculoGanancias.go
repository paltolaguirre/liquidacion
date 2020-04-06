package Ganancias

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-liquidacion/calculosAutomaticos"
	"github.com/jinzhu/copier"
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

const calculoremunerativos = -1
const calculoremunerativosmenosdescuentos = -3
const calculoremunerativosmasnoremunerativos = -4
const calculoremunerativosmasnoremunerativosmenosdescuentos = -5
const conceptoHorasExtrasCien = -6

const tipoconceptoremunerativos = -1
const tipoconceptodescuento = -3
const tipoconceptoretencion = -4

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

	cg.cloneRemplaceLiq()
	cg.obtenerLiquidacionesItemsPrimerQuincenaVacaciones()
	cg.recalcularImporteConceptosSiExisteHorasExtrasCien()
	cg.invocarCalculosLiquidacionAnual()
	calculo := (&CalculoRetencionDelMes{*cg}).getResult()
	return calculo

}

func (cg *CalculoGanancias) cloneRemplaceLiq() {
	var liquidacion structLiquidacion.Liquidacion
	copier.Copy(&liquidacion, &cg.Liquidacion)
	cg.Liquidacion = &liquidacion
	cg.copiarYReemplazarLiquidacionItems()
}

func (cg *CalculoGanancias) recalcularImporteConceptosSiExisteHorasExtrasCien() {
	if cg.existeHorasExtrasCien() {
		cg.recalcularImporteHorasExtrasCien()
		cg.recalcularImporteConceptos()
	}
}

func (cg *CalculoGanancias) copiarYReemplazarLiquidacionItems() {
	var arrayLiquidacionesItems []structLiquidacion.Liquidacionitem
	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		var liquidacionItem structLiquidacion.Liquidacionitem
		copier.Copy(&liquidacionItem, &cg.Liquidacion.Liquidacionitems[i])
		arrayLiquidacionesItems = append(arrayLiquidacionesItems, liquidacionItem)
	}

	cg.Liquidacion.Liquidacionitems = arrayLiquidacionesItems
}

func (cg *CalculoGanancias) obtenerLiquidacionesItemsPrimerQuincenaVacaciones() int {
	var liquidacionPrimerQuincena structLiquidacion.Liquidacion
	existeHorasExtrasCien := cg.existeHorasExtrasCien()
	items := len(cg.Liquidacion.Liquidacionitems)

	if cg.Liquidacion.Tipo.Codigo == "SEGUNDA_QUINCENA" || cg.Liquidacion.Tipo.Codigo == "MENSUAL" {
		mesliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
		anioLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()

		cg.Db.Set("gorm:auto_preload", true).Find(&liquidacionPrimerQuincena, "to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') = ? AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND id != ? AND legajoid = ? AND deleted_at is null", mesliquidacion, anioLiquidacion, strconv.Itoa(cg.Liquidacion.ID), cg.Liquidacion.Legajoid)

		for i := 0; i < len(liquidacionPrimerQuincena.Liquidacionitems); i++ {
			agregarLiquidacionItem := true
			liquidacionItem := liquidacionPrimerQuincena.Liquidacionitems[i]
			concepto := liquidacionItem.Concepto
			if existeHorasExtrasCien {
				if cg.esConceptoParaRecalcularImporte(concepto) {
					if cg.existeLiquidacionItemIntoArray(liquidacionItem) {
						agregarLiquidacionItem = false
					}
				}
			}
			if liquidacionItem.DeletedAt == nil && agregarLiquidacionItem {
				cg.Liquidacion.Liquidacionitems = append(cg.Liquidacion.Liquidacionitems, liquidacionItem)
			}
		}

	}

	return items
}

func (cg *CalculoGanancias) existeHorasExtrasCien() bool {
	var existeHorasExtrasCien = false
	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if concepto.ID == conceptoHorasExtrasCien {
			existeHorasExtrasCien = true
			break
		}
	}
	return existeHorasExtrasCien
}

func (cg *CalculoGanancias) existeLiquidacionItemIntoArray(liquidacionitem structLiquidacion.Liquidacionitem) bool {
	var noExisteLiquidacionItemIntoArray = false
	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		if cg.Liquidacion.Liquidacionitems[i].Concepto.ID == liquidacionitem.Concepto.ID {
			noExisteLiquidacionItemIntoArray = true
		}
	}
	return noExisteLiquidacionItemIntoArray
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
	(&CalculoImpuestoRetenido{*cg}).getResult()
	(&CalculoPagosACuenta{*cg}).getResult()
	(&CalculoSaldoAPagar{*cg}).getResult()
}

func (cg *CalculoGanancias) getfgSacCuotas(correspondeSemestre bool) float64 {
	var importeTotal, importeConcepto float64

	if correspondeSemestre {
		for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
			liquidacionitem := cg.Liquidacion.Liquidacionitems[i]

			if liquidacionitem.DeletedAt == nil {
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

					if *concepto.Tipoconceptoid == tipoconceptoretencion || *concepto.Tipoconceptoid == tipoconceptodescuento {
						importeConcepto = importeConcepto * -1
					}
					importeTotal = importeTotal + importeConcepto
				}
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

		if liquidacionitem.DeletedAt == nil {
			tipoconcepto := *liquidacionitem.Concepto.Tipoconceptoid
			importeconcepto := liquidacionitem.Importeunitario
			if importeconcepto != nil {

				if tipoconcepto == tipoconceptoremunerativos {
					totalRemunerativos = totalRemunerativos + *importeconcepto
				}
				if tipoconcepto == tipoconceptodescuento {
					totalDescuentos = totalDescuentos + *importeconcepto
				}
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

		if liquidacionitem.DeletedAt == nil {
			concepto := liquidacionitem.Concepto
			tipoimpuesto := obtenerTipoImpuesto(concepto, cg.Db)
			var mes float64 = 1

			if tipoimpuesto == tipoImpuestoALasGanancias && concepto.Codigo != "IMPUESTO_GANANCIAS" && concepto.Codigo != "IMPUESTO_GANANCIAS_DEVOLUCION" {
				if concepto.Prorrateo == true {
					mes = float64(cg.getfgMesesAProrratear(concepto))
				}

				if liquidacionitem.Importeunitario != nil {
					importeLiquidacionitem := *liquidacionitem.Importeunitario

					if *concepto.Tipoconceptoid == tipoconceptodescuento {
						importeLiquidacionitem = importeLiquidacionitem * -1
					}

					importeConcepto = importeLiquidacionitem / mes

				}
				importeTotal = importeTotal + importeConcepto

			}
		}

	}
	return importeTotal
}

func (cg *CalculoGanancias) obtenerImporteHorasExtrasCien() float64 {
	var importeConcepto, importeTotal float64

	for _, liquidacionItem := range cg.Liquidacion.Liquidacionitems {
		concepto := liquidacionItem.Concepto
		if concepto.ID == conceptoHorasExtrasCien {
			importeConcepto = *liquidacionItem.Importeunitario
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

		if mesdadoalta > mesperiodoliquidacion {
			importeTotal = 0
		} else {
			if mesdadobaja == 0 {
				importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
			} else {
				if mesdadobaja <= mesperiodoliquidacion {
					importeTotal = (valorfijo / 12) * float64(mesdadobaja-(mesdadoalta-1)) * (porcentaje / 100)
				} else {
					if mesdadobaja > mesperiodoliquidacion {
						importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
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
				importeTotal = (valorfijo / 12) * -float64(12-(mesdadoalta-1)) * (porcentaje / 100)
			} else {
				if mesdadobaja <= 12 {
					importeTotal = (valorfijo / 12) * float64(mesdadobaja-(mesdadoalta-1)) * (porcentaje / 100)
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

	var calculoGanancias CalculoGanancias

	for i := 0; i < len(liquidaciones); i++ {
		calculoGanancias.Liquidacion = &liquidaciones[i]
		cg.recalcularImporteConceptosSiExisteHorasExtrasCien()
		if calculoGanancias.existeHorasExtrasCien() {
			cg.recalcularImporteHorasExtrasCien()
			calculoGanancias.recalcularImporteConceptos()
			liquidaciones[i] = *calculoGanancias.Liquidacion
		}
	}

	return &liquidaciones
}

func (cg *CalculoGanancias) obtenerLiquidacionIgualAnioLegajoMesAnterior() *structLiquidacion.Liquidacion {
	var liquidacionMesAnterior structLiquidacion.Liquidacion
	liquidaciones := *cg.obtenerLiquidacionesIgualAnioLegajoMenorMes()
	if len(liquidaciones) > 0 {
		for _, liquidacion := range liquidaciones {
			if liquidacion.Tipo.Codigo == "MENSUAL" || liquidacion.Tipo.Codigo == "SEGUNDA_QUINCENA" {
				liquidacionMesAnterior = liquidacion
				break
			}
		}
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

func (cg *CalculoGanancias) getfgImporteTotalSiradigSegunTipoGrillaMesDesdeHasta(columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string) float64 {
	var importeTotal float64
	anioliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("01")

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + "* GREATEST((1 + LEAST(to_number(to_char(meshasta, 'MM'),'99'), " +
		mesLiquidacion + ") - to_number(to_char(mes, 'MM'),'99')), 0)) " +
		"FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = ts.siradigid WHERE stg.codigo = '" + tipodeducciondesgravacionsiradig + "' AND sdg.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anioliquidacion) + "' AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL AND ts.deleted_at IS NULL;"
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

func (cg *CalculoGanancias) retirarItemsPrimerQuincenaVacaciones(items int) {
	cg.Liquidacion.Liquidacionitems = cg.Liquidacion.Liquidacionitems[:items]
}

func (cg *CalculoGanancias) obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo(codigo string) float64 {
	var importeTotal float64

	anioLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion) - 1
	legajoID := *cg.Liquidacion.Legajoid

	sql := "select importe from acumulador where codigo = '" + codigo + "' and liquidacionitemid in (select id from liquidacionitem  where conceptoid = -29 and liquidacionid in (select ID from liquidacion where to_char(fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioLiquidacion) + "' and to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') = " + strconv.Itoa(mesLiquidacion) + " AND tipoid = -1 AND legajoid = " + strconv.Itoa(legajoID) + " order by fechaperiodoliquidacion))"
	cg.Db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (cg *CalculoGanancias) roundTo(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func (cg *CalculoGanancias) recalcularImporteHorasExtrasCien() {
	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if concepto.ID == conceptoHorasExtrasCien {
			importeConcepto := *liquidacionitem.Importeunitario / float64(2)
			cg.Liquidacion.Liquidacionitems[i].Importeunitario = &importeConcepto

		}
	}
}

func (cg *CalculoGanancias) recalcularImporteConceptos() {

	for i := 0; i < len(cg.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := cg.Liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if cg.esConceptoParaRecalcularImporte(concepto) {

			calculoAutomatico := calculosAutomaticos.NewCalculoAutomatico(concepto, cg.Liquidacion)
			calculoAutomatico.Hacercalculoautomatico()
			importeCalculadoConceptoID := cg.roundTo(calculoAutomatico.GetImporteCalculado(), 4)

			cg.Liquidacion.Liquidacionitems[i].Importeunitario = &importeCalculadoConceptoID

		}

	}
}

func (cg *CalculoGanancias) esConceptoParaRecalcularImporte(concepto *structConcepto.Concepto) bool {
	/*todo concepto con impuesto a las ganancias y tipo de calculo porcentaje (que utilice remunerativos) deberan recalcular su importe*/
	var esconceptopararecalcularimporte = false
	if concepto.Tipoimpuestogananciasid != nil {
		if concepto.Tipodecalculoid != nil {
			tipocalculo := *concepto.Tipodecalculoid
			if tipocalculo == calculoremunerativos || tipocalculo == calculoremunerativosmenosdescuentos || tipocalculo == calculoremunerativosmasnoremunerativos || tipocalculo == calculoremunerativosmasnoremunerativosmenosdescuentos {
				esconceptopararecalcularimporte = true
			}
		}
	}
	return esconceptopararecalcularimporte
}
