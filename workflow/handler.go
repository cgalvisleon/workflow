package workflow

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/utility"
)

var workFlows *WorkFlows

/**
* Load
* @return error
 */
func Load() error {
	if workFlows != nil {
		return nil
	}

	err := cache.Load()
	if err != nil {
		return err
	}

	err = event.Load()
	if err != nil {
		return err
	}

	workFlows = newWorkFlows()
	return nil
}

/**
* HealthCheck
* @return bool
**/
func HealthCheck() bool {
	if err := Load(); err != nil {
		return false
	}

	return workFlows.healthCheck()
}

/**
* New
* @param tag, version, name, description string, definition string, createdBy string
* @return *Flow
**/
func New(tag, version, name, description string, definition string, stop bool, createdBy string) *Flow {
	if err := Load(); err != nil {
		return nil
	}

	return workFlows.newFlowDefinition(tag, version, name, description, definition, stop, createdBy)
}

/**
* NewByFile
* @param tag, version, name, description string, filePath string, stop bool, createdBy string
* @return *Flow
**/
func NewByFile(tag, version, name, description string, filePath string, stop bool, createdBy string) *Flow {
	if err := Load(); err != nil {
		return nil
	}

	definition, err := os.ReadFile(filePath)
	if err != nil {
		logs.Error(err)
		definition = []byte("")
	}

	return workFlows.newFlowDefinition(tag, version, name, description, string(definition), stop, createdBy)
}

/**
* NewFn
* @param tag, version, name, description string, fn FnContext, createdBy string
* @return *Flow
**/
func NewFn(tag, version, name, description string, fn FnContext, stop bool, createdBy string) *Flow {
	if err := Load(); err != nil {
		return nil
	}

	return workFlows.newFlowFn(tag, version, name, description, fn, stop, createdBy)
}

/**
* LoadFlow
* @param params et.Json
* @return (*Flow, error)
**/
func LoadFlow(params et.Json) (*Flow, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	tag := params.Str("tag")
	name := params.Str("name")
	version := params.Str("version")
	if !utility.ValidStr(tag, 0, []string{""}) {
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "tag")
	}
	if !utility.ValidStr(name, 0, []string{""}) {
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "name")
	}
	if !utility.ValidStr(version, 0, []string{""}) {
		return nil, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "version")
	}

	description := params.Str("description")
	definition := params.Str("definition")
	stop := params.Bool("stop")
	createdBy := params.Str("createdBy")
	steps := params.ArrayJson("eteps")
	result := workFlows.newFlowDefinition(tag, version, name, description, definition, stop, createdBy)
	for i, step := range steps {
		name := step.Str("name")
		description := step.Str("description")
		if !utility.ValidStr(name, 0, []string{""}) {
			return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED_STEP, "name", i)
		}
		if !utility.ValidStr(description, 0, []string{""}) {
			return nil, fmt.Errorf(MSG_ATTRIBUTE_REQUIRED_STEP, "description", i)
		}

		definition := step.Str("definition")
		stop := step.Bool("stop")
		result.Step(name, description, definition, stop)
	}

	models := params.ArrayJson("models")
	for _, model := range models {
		dataBase := model.Str("database")
		name := model.Str("name")
		result.loadModel(dataBase, name)
	}

	return result, nil
}

/**
* LoadByTag
* @param tag string
* @return (*Flow, error)
**/
func LoadByTag(tag string) (*Flow, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	if getFlow == nil {
		return nil, fmt.Errorf(MSG_FLOW_NOT_FOUND)
	}

	result, err := getFlow(tag)
	if err != nil {
		return nil, err
	}

	workFlows.add(result)
	return result, nil
}

/**
* DeleteFlow
* @param tag string
* @return error
**/
func DeleteFlow(tag string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.deleteFlow(tag)
}

/**
* Run
* @param instanceId, tag string, startId int, tags et.Json, ctx et.Json, createdBy string
* @return et.Json, error
**/
func Run(instanceId, tag string, startId int, tags et.Json, ctx et.Json, createdBy string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.run(instanceId, tag, startId, tags, ctx, createdBy)
}

/**
* Continue
* @param instanceId string, tags et.Json, ctx et.Json, createdBy string
* @return et.Json, error
**/
func Continue(instanceId string, tags et.Json, ctx et.Json, createdBy string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	instance, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return et.Json{}, fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
	}

	return workFlows.run(instanceId, instance.Tag, instance.Current, tags, ctx, createdBy)
}

/**
* Reset
* @param instanceId string
* @return error
**/
func Reset(instanceId, createdBy string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.reset(instanceId, createdBy)
}

/**
* Rollback
* @param instanceId string
* @return et.Json, error
**/
func Rollback(instanceId string) (et.Json, error) {
	if err := Load(); err != nil {
		return et.Json{}, err
	}

	return workFlows.rollback(instanceId)
}

/**
* Stop
* @param instanceId string
* @return error
**/
func Stop(instanceId string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.stop(instanceId)
}

/**
* GetInstance
* @param instanceId string
* @return (*Instance, error)
**/
func GetInstance(instanceId string) (*Instance, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	result, exists := workFlows.loadInstance(instanceId)
	if !exists {
		return nil, fmt.Errorf(MSG_INSTANCE_NOT_FOUND)
	}

	return result, nil
}

/**
* DeleteInstance
* @param instanceId string
* @return error
**/
func DeleteInstance(instanceId string) error {
	if err := Load(); err != nil {
		return err
	}

	return workFlows.delete(instanceId)
}

/**
* HttpAll
* @params w http.ResponseWriter, r *http.Request
**/
func HttpAll(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := LoadFlow(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}
