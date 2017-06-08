/*
 * Este es el clásico ejemplo del cálculo del factorial
 * de un entero positivo n.
 *
 * (c) jsegovia@cnc.una.py
 */
var
   n : numerico
inicio
   imprimir ("\nCALCULO DE FACTORIAL",
             "\n--------------------",
             "\nIngrese un numero (0-n):")
   leer (n)
   si ( n >= 0 && n == int (n) ) {
      imprimir ("\n\n\n", n, "!=", fact (n))
   sino
      imprimir ("\nNo definido para ", n)
   }
fin


sub fact (n : numerico) retorna numerico
/*
 * Calcula el factorial de n. Imprime los valores que
 * se usaron en la etapa del cálculo.
 */
var
   r : numerico
inicio
   si ( n == 0 ) {
      r = 1
   sino
      imprimir ("\n", n, "! = ", n, " x (", n-1, "!)")
      r = n*fact(n-1)
   }
   retorna r
fin

sub fact2 (n : numerico) retorna numerico
/*
 * Una versión más compacta de fact(), que no imprime el
 * “rastro” de los valores intermedios.
 */
inicio
   retorna ifval ( n == 0,
                   1,
                   n*fact2(n-1))
fin
