package apiClientFormula

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/framework/configuracion"
	"io/ioutil"
	"net/http"
)

func ExecuteFormulaLiquidacion(authorization string, liquidacion *structLiquidacion.Liquidacion, formulaName string, concepto *structConcepto.Concepto) (float64, error) {

	config := configuracion.GetInstance()
	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioformula) + "formula/execute"

	executeBody := FormulaExecute{
		Context: Context{Currentliquidacion: *liquidacion, Currentconcepto: *concepto},
		Invoke: structFunction.Invoke{
			Functionname: formulaName,
		},
	}

	requestByte, err := json.Marshal(executeBody)

	if err != nil {
		return 0, errors.New("Error al convertir el body a string: " + err.Error())
	}

	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", url, requestReader)

	if err != nil {
		return 0, errors.New("Error al generar el request a " + url + ": " + err.Error())
	}

	req.Header.Add("Authorization", authorization)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, errors.New("Error al enviar el request a " + url + ": " + err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, errors.New("Error al leer la respuesta de " + url +": " + err.Error())
	}

	if res.StatusCode != http.StatusCreated {
		return 0, errors.New("No se pudo resolver la formula")
	}

	value := &structFunction.Value{}

	err = json.Unmarshal([]byte(string(body)), value)

	if err != nil {
		return 0, err
	}

	return value.Valuenumber, nil

}
