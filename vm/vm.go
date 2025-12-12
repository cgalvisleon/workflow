package vm

import (
	"github.com/cgalvisleon/et/et"
	"github.com/dop251/goja"
)

type Vm struct {
	*goja.Runtime
	Ctx et.Json
}

/**
* New
* Create a new vm
**/
func New() *Vm {
	result := &Vm{
		Runtime: goja.New(),
		Ctx:     et.Json{},
	}

	ToJson(result)
	ToString(result)
	Console(result)
	Fetch(result)
	Event(result)
	Cache(result)
	Model(result)
	Select(result)
	Query(result)
	return result
}

/**
* Run
* @param script string
* @return goja.Value, error
**/
func (v *Vm) Run(script string) (goja.Value, error) {
	if script == "" {
		return nil, nil
	}

	result, err := v.RunString(script)
	if err != nil {
		return nil, err
	}

	return result, nil
}
