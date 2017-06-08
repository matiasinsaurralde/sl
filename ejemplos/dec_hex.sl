/*
 * Lee un entero positivo e imprime el mismo número expresado
 * en base 16.
 *
 * (c) jsegovia@cnc.una.py
 */
var
   n = 0
   hex = “”
   k = 0
const
   DIG_HEX = “0123456789ABCDEF”
inicio
   imprimir (“Ingrese un entero positivo:”)
   leer (n)
   repetir
      hex = DIG_HEX [n % 16 + 1] + hex
      n = int (n/16)
   hasta ( n == 0 )
   imprimir (“\nHexadecimal=”, hex)
fin
