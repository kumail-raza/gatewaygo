package targets

import (
	"fmt"

	"github.com/minhajuddinkhan/gatewaygo/fhir"
)

var (
	//DefaultMapper DefaultMapper
	DefaultMapper = map[string]map[string]func(b []byte) (interface{}, error){

		"appointment": map[string]func(b []byte) (interface{}, error){
			"New": func(b []byte) (interface{}, error) {
				return fhir.NewAppointment(b)
			},
		},
		"patient": map[string]func(b []byte) (interface{}, error){
			"New": func(b []byte) (interface{}, error) {
				return fhir.NewFHIRPatient(b)
			},
		},
		"practitioner": map[string]func(b []byte) (interface{}, error){
			"New": func(b []byte) (interface{}, error) {
				return fhir.NewFHIRPractitioner(b)
			},
		},
		"encounter": map[string]func(b []byte) (interface{}, error){
			"New": func(b []byte) (interface{}, error) {
				return fhir.NewFHIREncounter(b)
			},
		},
	}
)

//DefaultTarget DefaultTarget
type DefaultTarget struct {
	DataModel string
	Event     string
	ToFHIR    func(b []byte) (interface{}, error)
}

//NewDefaultTarget NewDefaultTarget
func NewDefaultTarget(dataModel, event string) DefaultTarget {

	d := DefaultTarget{
		DataModel: dataModel,
		Event:     event,
		ToFHIR: func(b []byte) (interface{}, error) {

			var fhir interface{}
			if fn, ok := DefaultMapper[dataModel][event]; ok {
				result, err := fn(b)
				if err != nil {
					return fhir, fmt.Errorf("Cannot Map for Default Target. Error: %s", err.Error())
				}
				return result, nil
			}
			return fhir, fmt.Errorf("Mapper not configured for Default Mapper {DataModel: %s, Event: %s}", dataModel, event)

		},
	}
	return d
}
