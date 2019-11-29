package main

import (
	"sort"
	"strings"
)

const (
	comprimentoMapaFreq = 26 // Letras maiúsculas.
)

type (
	// MapaFrequencia mantém uma letra no mapa de frequência
	// Somente caracteres 'comprimentoMapaFreq' são suportados.
	MapaFrequencia [comprimentoMapaFreq]int
)

// candidatos lê uma slice de palavras e produz uma lista de palavras candidatas
// (ou seja, palavras que podem ser anagramas para nossa frase).
func candidatos(palavras []string, frase string, tamMinPalavras, tamMaxPalavras int) []string {
	var cand []string
	frasemap := mapaFreq(&frase)
	tamfrase := len(frase)

loopPalavras:
	for _, p := range palavras {
		tampalavra := len(p)

		// Próxima palavra imediatamente se a palavra for maior que a frase.
		if tampalavra > tamfrase {
			continue
		}
		// Rejeitar a palavra se estiver fora dos limites desejados
		if tampalavra < tamMinPalavras || tampalavra > tamMaxPalavras {
			continue
		}
		// Ignore qualquer coisa que não esteja em [A-Z].
		p = strings.ToUpper(p)
		for _, r := range p {
			if r < 'A' || r > 'Z' {
				continue loopPalavras
			}
		}
		if contemNoMapa(&frasemap, &p) {
			cand = append(cand, p)
		}
	}

	sort.Sort(peloTam(cand))
	return cand
}

// O mapaFreq cria um mapa de frequência de todas as letras da palavra.
// Ele assume apenas letras maiúsculas como entrada e
// usa uma slice em vez de maps por motivos de desempenho.
func mapaFreq(palavra *string) MapaFrequencia {
	var m MapaFrequencia

	for ix := 0; ix < len(*palavra); ix++ {
		r := (*palavra)[ix]
		idx := r - 'A'
		m[idx]++
	}
	return m
}

// tamMapa retorna o comprimento do mapa.
func tamMapa(m MapaFrequencia) int {
	var tam int
	for i := 0; i < comprimentoMapaFreq; i++ {
		tam += m[i]
	}
	return tam
}

// contemNoMapa retorna true se o mapa 'a' contiver a string.
func contemNoMapa(a *MapaFrequencia, palavra *string) bool {
	var smap MapaFrequencia

	for i := 0; i < len(*palavra); i++ {
		idx := (*palavra)[i] - 'A'
		smap[idx]++
	}
	for i := 0; i < comprimentoMapaFreq; i++ {
		if smap[i] > (*a)[i] {
			return false
		}
	}
	return true
}

// subtrairMapa retorna um mapa representando o mapa a - map b.
func subtrairMapa(m MapaFrequencia, palavras []string) MapaFrequencia {
	total := MapaFrequencia{}

	for i := 0; i < len(palavras); i++ {
		for j := 0; j < len(palavras[i]); j++ {
			idx := palavras[i][j] - 'A'
			total[idx]++
		}
	}
	for i := 0; i < comprimentoMapaFreq; i++ {
		total[i] = m[i] - total[i]
	}
	return total
}

// mapaEstaVazio retorna true se o mapa estiver vazio, caso contrário retorna false.
func mapaEstaVazio(m MapaFrequencia) bool {
	for i := 0; i < comprimentoMapaFreq; i++ {
		if m[i] > 0 {
			return false
		}
	}
	return true
}

// anagramas gera recursivamente uma lista de anagramas para a lista especificada de
// candidatos, começando com 'base' como raiz. Se depht 'profundidade' for especificada,
// a recursão será interrompida neste nível. Isso limita essencialmente o número de
// palavras em um anagrama. Esta função pode levar um tempo muito longo se o
// o número de palavras candidatas for muito grande.
func anagramas(frasemap MapaFrequencia, cand []string, base []string, numpalavras, maxpalavras int) []string {
	var ret []string

	// profundidade máxima de recursão (número de palavras)
	if numpalavras > maxpalavras {
		return nil
	}
	numpalavras++

	// maparestante é o que sobrou na palavra depois que voce "subtraiu" o mapa de outra.
	// Tipo: BANANA - ANA = BAN
	maparestante := subtrairMapa(frasemap, base)

	// Combinação perfeita.
	if mapaEstaVazio(maparestante) {
		return append(ret, strings.Join(base, " "))
	}

	charsRestantes := tamMapa(maparestante)

	for ix := 0; ix < len(cand); ix++ {
		palavraAtual := cand[ix]
		// A lista de entrada de palavras é classificada pelo comprimento da palavra.
		// Se o comprimento da base atual + a palavra atual excede o comprimento total
		// da frase, não existem mais anagramas a partir deste ponto.
		if len(palavraAtual) > charsRestantes {
			break
		}

		// Recorra apenas se palavraAtual se encaixar nos caracteres restantes.
		if !contemNoMapa(&maparestante, &palavraAtual) {
			continue
		}

		// Nova base é a nossa base atual + nova palavra.
		novabase := append(base, palavraAtual)
		r := anagramas(frasemap, cand[ix+1:], novabase, numpalavras, maxpalavras)
		if r != nil {
			ret = append(ret, r...)
		}
	}
	return ret
}
