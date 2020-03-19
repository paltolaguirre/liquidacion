package apiClientFormula

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type FormulaExecute struct {
	structFunction.Invoke
	Context Context `json:"context"`
}

type Executor struct {
	db      *gorm.DB
	stack   [][]structFunction.Value
	context *Context //[]byte
}

type Context struct {
	Currentliquidacion structLiquidacion.Liquidacion `json:"currentliquidacion"`
}
