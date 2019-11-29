package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"sort"
	"strings"
)

var (
	optCandidates bool
	optCPUProfile string
	optDict       string
	optMinWordLen int
	optMaxWordLen int
	optMaxWordNum int
	optSilent     bool
	optSortLines  bool
	optSortWords  bool
)

// Sanitize converte a string de entrada em maiúsculas e remove todos os caracteres
// que não corresponde a [A-Z].
func sanitize(s string) (string, error) {
	re, err := regexp.Compile("[^A-Z]")
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(strings.ToUpper(s), ""), nil
}

// sortWords lê uma slice de strings e classifica cada linha por palavra.
func sortWords(lines []string) {
	for idx, line := range lines {
		w := strings.Split(line, " ")
		sort.Strings(w)
		lines[idx] = strings.Join(w, " ")
	}
}

// printAnagrams imprime todos os anagramas usando as opções de (sort) classificação da linha de comando.
func printAnagrams(an []string, sortlines, sortwords bool) {
	if optSortWords {
		sortWords(an)
	}
	if optSortLines {
		sort.Strings(an)
	}
	for _, w := range an {
		fmt.Println(w)
	}
}

// readDict lê o dicionário na memória (uma palavra por linha) e
// retorna uma slice de strings com as palavras.
func readDict(dfile string) ([]string, error) {
	// lê o arquivo inteiro na memória.
	buf, err := ioutil.ReadFile(dfile)
	if err != nil {
		return nil, err
	}

	// Split - divida a entrada em novas linhas e gere a lista de palavras candidatas.
	words := strings.Split(strings.TrimRight(string(buf), "\n"), "\n")
	return words, nil
}

// printCandidates imprime as palavras candidatas (e alternativas).
func printCandidates(cand []string) {
	for _, w := range cand {
		fmt.Println(w)
	}
}

func main() {
	log.SetFlags(0)

	flag.BoolVar(&optCandidates, "candidates", false, "just show candidate words (don't anagram)")
	flag.StringVar(&optCPUProfile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&optDict, "dict", "words.txt", "dictionary file")
	flag.IntVar(&optMinWordLen, "minlen", 1, "minimum word length")
	flag.IntVar(&optMaxWordLen, "maxlen", 64, "maximum word length")
	flag.IntVar(&optMaxWordNum, "maxwords", 16, "maximum number of words (0=no maximum)")
	flag.BoolVar(&optSilent, "silent", false, "don't print results.")
	flag.BoolVar(&optSortLines, "sortlines", false, "(also) sort the output by lines")
	flag.BoolVar(&optSortWords, "sortwords", true, "(also) sort the output by words")

	// Uso personalizado.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use: %s [flags] expression_to_anagram\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}
	// Cria perfil de CPU
	if optCPUProfile != "" {
		f, err := os.Create(optCPUProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	phrase, err := sanitize(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	words, err := readDict(optDict)
	if err != nil {
		log.Fatal(err)
	}

	// Gere uma lista de palavras candidatas e alternativas.
	cand := candidatos(words, phrase, optMinWordLen, optMaxWordLen)

	if optCandidates {
		printCandidates(cand)
		os.Exit(0)
	}

	// Anagrama e impressão classificados por palavra (e opcionalmente, por linha).
	var an []string
	an = anagramas(mapaFreq(&phrase), cand, an, 0, optMaxWordNum)

	if !optSilent {
		printAnagrams(an, optSortLines, optSortWords)
	}
}
