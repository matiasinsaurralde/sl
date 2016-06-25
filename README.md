# sl

Intérprete experimental de SL, basado en Go.

## Ideas

* Agregar tests y lograr que corra la mayoría de los ejemplos del libro.
* Exponer algunas cosas de Go o incrustar C (mediante cgo).

## Uso

Programa de prueba:

```
% cat ejemplos/holamundo.sl

/*
 * Comentario
 * De prueba
 */

var
a : numerico
b : numerico
c : numerico
x=0

ll = 1

zz=1+2

inicio
  imprimir("hola")

  imprimir("mundo")
fin
```

Salida:

```
% go run test.go
AST Test Program

Parse...

- Found a program:  hola

- Found a comment...
  Comentario
  De prueba

- Declaration:

 * Node: &{a [] 0 0}
 * Node: &{b [] 0 0}
 * Node: &{c [] 0 0}
 * Node: &{x [] 0 0}
 * Node: &{ll [] 0 0}
 * Node: &{zz [] 0 0}

- Found a block

 * Statement:   imprimir("hola")
 * Statement:   imprimir("mundo")

- Scope: &{[0xc8200142c0 0xc820014340 0xc8200143c0 0xc820014480 0xc820014500 0xc820014580] [0xc820010440]}


Ast.File:
&{ejemplos/holamundo.sl hola 0xc820032028 0xc820012300 [{
  Comentario
  De prueba
  15 47}]}

Running...

Declaring: &{a [] 0 0}
Declaring: &{b [] 0 0}
Declaring: &{c [] 0 0}
Declaring: &{x [] 0 0}
Declaring: &{ll [] 0 0}
Declaring: &{zz [] 0 0}

Main...

vars map[]
Evaluate: &{imprimir [0xc8200125d0] 0 0}

Call: imprimir (function)

hola
Evaluate: &{imprimir [0xc820012600] 0 0}

Call: imprimir (function)

mundo
%
```

# Licencia

[MIT](https://github.com/matiasinsaurralde/sl/blob/master/LICENSE)
