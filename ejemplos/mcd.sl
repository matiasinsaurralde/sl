/*
 * Cálculo del máximo común divisor utilizando el algoritmo
 * de Euclides.
 *
 * (c) jsegovia@cnc.una.py
 */
var
   a, b : numerico
inicio
   imprimir ("Ingrese dos enteros positivos:")
   leer (a, b)
   si ( (a < 1) or (b < 1) ) {
       terminar (“\nLos valores ingresados deben ser positivos”)
   }
   mientras (a <> b ) {
      si ( a > b ) {
         a = a - b
      sino
         b = b - a
      }
   }
   imprimir (“\nEl MCD es “, a)
fin
