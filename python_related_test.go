package main

import (
	py3 "github.com/DataDog/go-python3"
	"testing"
)

func TestPythonRunning(t *testing.T) {
	pySetup()
	py3.PyRun_SimpleString("import sys\nsys.testA = True")
	if py3.PySys_GetObject("testA") != py3.Py_True {
		t.Error("Cannot set system values.")
	}
	pyStop()
}
