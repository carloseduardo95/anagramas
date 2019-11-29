// Unit tests para anagramarama.

package main

import (
	"sort"
	"testing"
)

func TestAnagram(t *testing.T) {
	casetests := []struct {
		phrase     string
		dictFile   string
		wantFile   string
		minWordLen int
		maxWordLen int
		maxWordNum int
		wantError  bool
	}{
		// Meu Teste.
		{
			phrase:     "carro",
			dictFile:   "testdata/words.txt",
			wantFile:   "testdata/meuresults.txt",
			minWordLen: 1,
			maxWordLen: 64,
			maxWordNum: 16,
		},
		// Uma thread.
		{
			phrase:     "marco paganini ab",
			dictFile:   "testdata/words.txt",
			wantFile:   "testdata/results.txt",
			minWordLen: 1,
			maxWordLen: 64,
			maxWordNum: 16,
		},
		// Multiplas threads.
		{
			phrase:     "marco paganini ab",
			dictFile:   "testdata/words.txt",
			wantFile:   "testdata/results.txt",
			minWordLen: 1,
			maxWordLen: 64,
			maxWordNum: 16,
		},
		// Conjunto mínimo e máximo de comprimento de palavra.
		{
			phrase:     "marco paganini ab",
			dictFile:   "testdata/words.txt",
			wantFile:   "testdata/results-min4-max5.txt",
			minWordLen: 4,
			maxWordLen: 5,
			maxWordNum: 16,
		},
		// Limita o número de palavras a 3.
		{
			phrase:     "lorem ipsum dolor sit",
			dictFile:   "testdata/words.txt",
			wantFile:   "testdata/results-3words.txt",
			minWordLen: 1,
			maxWordLen: 64,
			maxWordNum: 3,
		},
		// Nome de dicionário inválido (erro)
		{
			phrase:     "marco paganini ab",
			dictFile:   "INVALIDFILE",
			wantFile:   "testdata/results.txt",
			minWordLen: 1,
			maxWordLen: 64,
			maxWordNum: 16,
			wantError:  true,
		},
	}

	for _, tt := range casetests {
		phrase, err := sanitize(tt.phrase)
		if err != nil {
			t.Fatalf("error sanitizing phrase: %v", err)
		}

		// Arquivo de resultados
		want, err := readDict(tt.wantFile)
		if err != nil {
			t.Fatalf("error reading results file: %v", err)
		}

		words, err := readDict(tt.dictFile)
		if !tt.wantError {
			if err != nil {
				t.Fatalf("Got error %q want no error", err)
			}

			// Gere uma lista de palavras candidatas e alternativas.
			cand := candidatos(words, phrase, tt.minWordLen, tt.maxWordLen)
			got := anagramas(mapaFreq(&phrase), cand, []string{}, 0, tt.maxWordNum)

			lenGot := len(got)
			lenWant := len(want)
			if lenGot != lenWant {
				t.Fatalf("anagram lists have different lengths. Phrase %q, Got %d lines, want %d lines.", tt.phrase, lenGot, lenWant)
			}

			// Classifique a saída por palavras e depois por linha.
			sortWords(got)
			sort.Strings(got)

			for ix := 0; ix < lenGot; ix++ {
				if got[ix] != want[ix] {
					t.Fatalf("diff: phrase %q, line %d, Got %q, want %q.", tt.phrase, ix, got[ix], want[ix])
				}
			}
			continue
		}

		// Aqui, queremos ver um erro.
		if err == nil {
			t.Errorf("Got no error, want error")
		}
	}
}
