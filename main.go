package main

import (
	"fmt"
)

type Veiculo interface {
	Andar()
}

type Carro struct {
	Modelo string
	Ano    int
}

func (c Carro) Andar() {
	fmt.Println("O carro está andando")
}

func VemAndarComigo(v Veiculo) {
	v.Andar()
}

func main() {

	carro1 := Carro{Modelo: "Fusca", Ano: 1970}

	VemAndarComigo(carro1)
}
