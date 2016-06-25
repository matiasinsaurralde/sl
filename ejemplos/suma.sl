/*
  Calcular e imprimir la suma de los números pares comprendidos entre 1 y un n dado.

  Las variables utilizadas son:
    n : para almacenar el valor proveido por el usuario como tope.
    suma : donde almacenamos el total calculado.
    k : contador auxiliar para el ciclo.

  AUTOR : Juan Segovia (jsegovia@cnc.una.py)
*/

var
  n, suma, k : numerico

inicio
  imprimir (“\nSuma de numeros pares entre 1 y n.\nPor favor ”,
    “ingrese un valor para n: “)

  leer (n)
  suma = 0
  desde k=2 hasta n paso 2
  {
    suma = suma + k
  }
  imprimir ("\nLa suma es ", suma)
fin
