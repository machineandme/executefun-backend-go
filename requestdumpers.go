package main

import py3 "github.com/DataDog/go-python3"

func maybeArrayMapAsPyDict(items map[string][]string) *py3.PyObject {
	mapAsDict := py3.PyDict_New()
	for key, values := range items {
		if len(values) > 1 {
			args := py3.PyTuple_New(len(values))
			for i, v := range values {
				py3.PyTuple_SetItem(args, i, py3.PyUnicode_FromString(v))
			}
			mapAsDict.SetItem(py3.PyUnicode_FromString(key), args)
		} else {
			mapAsDict.SetItem(
				py3.PyUnicode_FromString(key),
				py3.PyUnicode_FromString(values[0]),
			)
		}
	}
	return mapAsDict
}

func mapAsPyDict(items map[string]string) *py3.PyObject {
	mapAsDict := py3.PyDict_New()
	for key, value := range items {
		mapAsDict.SetItem(
			py3.PyUnicode_FromString(key),
			py3.PyUnicode_FromString(value),
		)
	}
	return mapAsDict
}

func pyMapToGo(userData *py3.PyObject) (resultingMap map[string]string) {
	resultingMap = make(map[string]string)
	var key, val *py3.PyObject
	pos := 0
	for py3.PyDict_Next(userData, &pos, &key, &val) {
		realKey := py3.PyUnicode_AsUTF8(key.Str())
		realVal := py3.PyUnicode_AsUTF8(val.Str())
		resultingMap[realKey] = realVal
	}
	return
}
