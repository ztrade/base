package common

import (
	"testing"
)

func TestParam(t *testing.T) {
	var str1, str2 string
	var int1, int2 int
	var float1, float2 float64
	params := []Param{
		StringParam("str1", "str test", "just a simple string", "a", &str1),
		StringParam("str2", "str test", "enum string", "B", &str2,
			Entry{Label: "A", Value: "A"},
			Entry{Label: "B", Value: "B"}),
		IntParam("int1", "int1 test", "just a simple int", 1, &int1),
		IntParam("int2", "int2 test", "enum int", 1, &int2,
			Entry{Label: "A", Value: 1},
			Entry{Label: "B", Value: 2}),
		FloatParam("float1", "float1 test", "just a simple int", 1, &float1),
		FloatParam("float2", "float2 test", "enum float", 1, &float2,
			Entry{Label: "A", Value: 1.0},
			Entry{Label: "B", Value: 2.0}),
	}

	str := `{"str1": "str1", "str2":"A", "int1": 10, "int2": 1, "float1": 3, "float2": 2.0}`
	rets, err := ParseParams(str, params)
	if err != nil {
		t.Fatal(err.Error())
	}
	if str1 != "str1" || str2 != "A" || int1 != 10 || int2 != 1 || float1 != 3 || float2 != 2.0 {
		t.Fatal("value not match", str1, str2, int1, int2, float1, float2)
	}
	if rets["str1"] != str1 || rets["str2"] != str2 || rets["int1"] != int1 || rets["int2"] != int2 || rets["float1"] != float1 || rets["float2"] != float2 {
		t.Fatal("value not match", str1, str2, int1, int2, float1, float2, rets)
	}
}
