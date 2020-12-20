package main

import py3 "github.com/DataDog/go-python3"

func callSnake(someDict *py3.PyObject) *py3.PyObject {
	defer someDict.DecRef()
	py3.PySys_SetObject("scope", someDict)
	py3.PyRun_SimpleString("sys.response = (handler(sys.scope))")
	return py3.PySys_GetObject("response")
}
