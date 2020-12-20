package main

import py3 "github.com/DataDog/go-python3"

func SerializePyObj(target *py3.PyObject) string {
	return dumpValue(target)
}

func dumpValue(val *py3.PyObject) string {
	var realVal string
	switch val.Type() {
	case py3.Dict:
		realVal = serPyDict(val)
	case py3.Float, py3.Long:
		realVal = py3.PyUnicode_AsUTF8(val.Str())
	case py3.List, py3.Tuple:
		realVal = serPyList(val)
	case py3.Py_False:
		realVal = "false"
	case py3.Py_True:
		realVal = "true"
	case py3.Py_None:
		realVal = "null"
	default:
		realVal = "\"" + py3.PyUnicode_AsUTF8(val.Str()) + "\""
	}
	return realVal
}

func serPyList(target *py3.PyObject) string {
	result := "["
	if target.Type() == py3.List {
		target = py3.PyList_AsTuple(target)
	}
	size := py3.PyTuple_Size(target)
	if size == 0 {
		return "[]"
	}
	for i := 0; i < size; i++ {
		if i != 0 {
			result += ","
		}
		result += dumpValue(py3.PyTuple_GetItem(target, i))
	}
	return result + "]"
}

func serPyDict(target *py3.PyObject) string {
	result := "{"
	var key, val *py3.PyObject
	pos := 0
	for py3.PyDict_Next(target, &pos, &key, &val) {
		if pos != 1 {
			result += ","
		}
		realKey := py3.PyUnicode_AsUTF8(key)
		realVal := dumpValue(val)
		result += "\"" + realKey + "\":" + realVal
	}
	if pos == 0 {
		return "{}"
	}
	return result + "}"
}
