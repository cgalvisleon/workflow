package main

import (
	"net/http"

	"github.com/cgalvisleon/et/claim"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/workflow/workflow"
)

func main() {
	flowDefinition()
	// utility.AppWait()

	logs.Logf("test", "Fin de flow")

}

func flowDefinition() {
	workflow.NewByFile("report:set", "1.0.0", "Registro de reportes", "", "./step0.js", false, "test").
		StepByFile("Step 1", "Step 1", "./step1.js", true).
		AddModel("postgres", "report")

	result, err := workflow.Run("1234", "report:set", 0, et.Json{}, et.Json{
		"id":   "1234",
		"name": "report",
	}, "createdBy")
	if err != nil {
		logs.Error(err)
	} else {
		logs.Logf("Result", result.ToString())
	}
}

func HttpVenta(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	tag := r.PathValue("tag")
	serviceId := body.Str("serviceId")
	tags := et.Json{
		"cedula": "91499023",
		"codigo": "112342",
	}
	step := body.Int("step")
	createdBy := claim.ClientName(r)
	result, err := workflow.Run(serviceId, tag, step, tags, body, createdBy)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok:     true,
		Result: result,
	})

}
