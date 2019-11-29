// Sort interfaces e funções.
package main

// peloTam define um tipo para classificar(sort) uma slice de strings
// pelo tamanho de cada elemento.
type peloTam []string

func (x peloTam) Len() int {
	return len(x)
}

func (x peloTam) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x peloTam) Less(i, j int) bool {
	tami := len(x[i])
	tamj := len(x[j])
	return tami < tamj
}
