/*
 * Cálculo del mínimo común múltiplo, usando la relación
 *
 *              a * b
 * MCM(a,b) =  ------------
 *            MCD (a, b)
 *
 * (c) jsegovia@cnc.una.py
 */
var
   a, b : numerico
inicio
   imprimir ("Ingrese dos enteros positivos:")
   leer (a, b)
   imprimir ("\nEl MCM de ", a, " y ", b, " es ",
             (a*b) / MCD (a, b))
fin


sub MCD (a, b : numerico) retorna numerico
inicio
   mientras (a <> b ) {
      si ( a > b ) {
         a = a - b
      sino
         b = b - a
      }
   }
   retorna a
fin
