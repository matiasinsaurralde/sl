/*
 * Dado un arreglo bidimensional (matriz), genera e imprime
 * su traspuesta.
 *
 * (c) jsegovia@cnc.una.py
 */
var
   /*
    * Las 3 últimas filas de M serán iguales.
    */
   M : matriz [5, 3] = {{7, 12, 5},
                        {1, 4, 22},
                        {6, 20, 13},
                        ...
                       }
   T : matriz [*,*] numerico
inicio
   impr_mat (“Matriz original:\n”, M)
   transponer (M, T)
   impr_mat (“\nLa traspuesta es:\n”, T)
fin


sub transponer ( M : matriz [*,*] numerico
             ref R : matriz [*,*] numerico)
/*
 * trasponer() produce la transpuesta de M y lo deposita
 * en R. 
 * M puede tener cualquier tamaño, con tal de que
 * sea bidimensional y rectangular (cantidad igual de
 * elementos por cada fila).
 * R debe ser un arreglo abierto.
 */
var
   cant_filas = alen (M)
   cant_cols  = alen (M [1])
   filas, cols : numerico
inicio
   /*
    * Nótese que las filas y columnas están en orden
    * inverso en el siguiente dim().
    */
   dim (R, cant_cols, cant_filas)
   desde filas=1 hasta cant_filas {
      desde cols=1 hasta cant_cols {
         R [filas, cols] = M [cols, filas]
      }
   }
fin


sub impr_mat (msg : cadena; M : matriz [*,*] numerico)
var
   k = 0
inicio
   imprimir (msg)
   desde k=1 hasta alen (M) {
      imprimir (M [k], “\n”)
   }
fin
