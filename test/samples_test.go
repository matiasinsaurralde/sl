package test

import (
	"testing"
)

func TestHolamundo(t *testing.T) {
	result, err := RunSLProgram("../samples/holamundo.sl", "")
	if err != nil {
		t.Fatalf("failed to run holamundo.sl: %v", err)
	}

	expected := "hola\nmundo\n"
	AssertOutput(t, result.Stdout, expected, "holamundo.sl")
}

func TestSuma(t *testing.T) {
	result, err := RunSLProgram("../samples/suma.sl", "10\n")
	if err != nil {
		t.Fatalf("failed to run suma.sl: %v", err)
	}

	expected := "\nSuma de numeros pares entre 1 y n.\nPor favor ingrese un valor para n: \nIngrese valor para n: \nLa suma es 30\n"
	AssertOutput(t, result.Stdout, expected, "suma.sl")
}

func TestSumaWithDifferentInput(t *testing.T) {
	result, err := RunSLProgram("../samples/suma.sl", "6\n")
	if err != nil {
		t.Fatalf("failed to run suma.sl with input 6: %v", err)
	}

	expected := "\nSuma de numeros pares entre 1 y n.\nPor favor ingrese un valor para n: \nIngrese valor para n: \nLa suma es 12\n"
	AssertOutput(t, result.Stdout, expected, "suma.sl with input 6")
}

func TestSumaWithZeroInput(t *testing.T) {
	result, err := RunSLProgram("../samples/suma.sl", "0\n")
	if err != nil {
		t.Fatalf("failed to run suma.sl with input 0: %v", err)
	}

	expected := "\nSuma de numeros pares entre 1 y n.\nPor favor ingrese un valor para n: \nIngrese valor para n: \nLa suma es 0\n"
	AssertOutput(t, result.Stdout, expected, "suma.sl with input 0")
}
