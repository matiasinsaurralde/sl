/*
 * Lee un vector de 10 elementos numéricos, ordena
 * ascendentemente sus elementos y los imprime.
 * La ordenación se hace con el algoritmo de la “burbuja”.
 *
 * (c) jsegovia@cnc.una.py
 */
var
   A : vector [10] numerico
   m, n : numerico
inicio
   imprimir (“\nIngrese “, alen (A), “ números separados por comas:\n”)
   leer (A)
   desde m=1 hasta alen(A)-1 {
      desde n=m+1 hasta alen (A) {
         si ( A [m] > A [n] ) {
            intercambiar (A [m], A [n])
         }
      }
   }
   imprimir ("\nEl vector ordenado es:\n”, A)
fin
