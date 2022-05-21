package common

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/ztrade/trademodel"
)

type CandleFn func(candle *trademodel.Candle)

type Entry struct {
	Value interface{}
	Label string
}

type Param struct {
	Name     string
	Type     string
	Label    string
	Info     string
	DefValue interface{}
	Enums    []Entry
	ptr      interface{}
}

func StringParam(name, label, info, defValue string, ptr *string) Param {
	*ptr = defValue
	return Param{Name: name, Type: "string", Label: label, Info: info, DefValue: defValue, ptr: ptr}
}

func StringEnumParam(name, label, info, defValue string, enums []Entry, ptr *string) Param {
	*ptr = defValue
	return Param{Name: name, Type: "string", Label: label, Info: info, DefValue: defValue, Enums: enums, ptr: ptr}
}

func IntParam(name, label, info string, defValue int, ptr *int) Param {
	*ptr = defValue
	return Param{Name: name, Type: "int", Label: label, Info: info, DefValue: defValue, ptr: ptr}
}

func IntEnumParam(name, label, info string, defValue int, enums []Entry, ptr *int) Param {
	*ptr = defValue
	return Param{Name: name, Type: "int", Label: label, Info: info, DefValue: defValue, Enums: enums, ptr: ptr}
}

func FloatParam(name, label, info string, defValue float64, ptr *float64) Param {
	*ptr = defValue
	return Param{Name: name, Type: "float", Label: label, Info: info, DefValue: defValue, ptr: ptr}
}

func FloatEnumParam(name, label, info string, defValue float64, enums []Entry, ptr *float64) Param {
	*ptr = defValue
	return Param{Name: name, Type: "float", Label: label, Info: info, DefValue: defValue, Enums: enums, ptr: ptr}
}

func BoolParam(name, label, info string, defValue bool, ptr *bool) Param {
	*ptr = defValue
	return Param{Name: name, Type: "bool", Label: label, Info: info, DefValue: defValue, ptr: ptr}
}

func BoolEnumParam(name, label, info string, defValue bool, enums []Entry, ptr *bool) Param {
	*ptr = defValue
	return Param{Name: name, Type: "bool", Label: label, Info: info, DefValue: defValue, Enums: enums, ptr: ptr}
}

func ParseParams(str string, params []Param) (data ParamData, err error) {
	data = make(ParamData)
	sj := simplejson.New()
	err = sj.UnmarshalJSON([]byte(str))
	if err != nil {
		return
	}

	var temp *simplejson.Json
	var ok, boolV bool
	var strV string
	var intV int
	var floatV float64
	for _, v := range params {
		if v.ptr == nil {
			data[v.Name] = sj.Get(v.Name).Interface()
			return
		}
		temp, ok = sj.CheckGet(v.Name)
		if !ok {
			continue
		}

		switch ptr := v.ptr.(type) {
		case *string:
			strV, err = temp.String()
			if err != nil {
				return
			}
			*ptr = strV
			data[v.Name] = strV
		case *float64:
			floatV, err = temp.Float64()
			if err != nil {
				return
			}
			*ptr = floatV
			data[v.Name] = floatV
		case *int:
			intV, err = temp.Int()
			if err != nil {
				return
			}
			*ptr = intV
			data[v.Name] = intV
		case *bool:
			boolV, err = temp.Bool()
			if err != nil {
				return
			}
			*ptr = boolV
			data[v.Name] = boolV
		default:
			err = fmt.Errorf("unsupport value type: %##v", ptr)
			return
		}
	}
	return
}

type ParamData map[string]interface{}

func (d ParamData) GetString(key, defaultValue string) string {
	v, ok := d[key]
	if !ok {
		return defaultValue
	}
	ret := v.(string)
	if ret == "" {
		return defaultValue
	}
	return ret
}
func (d ParamData) GetFloat(key string, defaultValue float64) float64 {
	v, ok := d[key]
	if !ok {
		return defaultValue
	}
	ret := v.(float64)
	if ret == 0 {
		return defaultValue
	}
	return ret
}
