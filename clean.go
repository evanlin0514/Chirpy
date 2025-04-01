package main

import(
	"strings"
)

type cleanBody struct {
	Body string `json:"cleaned_body"`
}
type parameters struct {
	Body string `json:"body"`
}

func cleanInput(str string) cleanBody{
	clean := cleanBody{
		Body: str,
	}

	banMaps := make(map[string]bool)
	banWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, ban := range banWords{
		banMaps[ban] = true
	}

	words := strings.Split(clean.Body, " ")
	for i, word := range words {
		if banMaps[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	clean.Body = strings.Join(words, " ")
	return clean
}