package marisfrolg_utils

import (
	"reflect"
)

func GetModelProperty(model interface{}, propertyName string) (result interface{}) {
	types := reflect.TypeOf(model)
	vals := reflect.ValueOf(model)
	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		if field.Name == propertyName {
			result = vals.Field(i).Interface()
			break
		}
	}
	return
}
