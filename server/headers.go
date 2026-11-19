package server

import (
	"net/http"
	"sort"
)

// getHeaders return a response with all the request headers
func (s *Server) getHeaders() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		headers := make(map[string][]string)
		for name, values := range r.Header {
			headers[name] = values
		}
		s.Respoder.SendOK(w, headers)
		return
	}
}

// sortHeadersAlphaBetically get headers from request and arrage it alphabetically
func (s *Server) sortHeadersAlphaBetically() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		alphabeticallyHeadersKeys := []string{}
		alphabeticallyHeadersValues := []string{}
		for name, values := range r.Header {
			for _, value := range values {
				alphabeticallyHeadersKeys = append(alphabeticallyHeadersKeys, name)
				alphabeticallyHeadersValues = append(alphabeticallyHeadersValues, value)
			}
		}
		sort.Strings(alphabeticallyHeadersKeys)
		sort.Strings(alphabeticallyHeadersValues)
		s.Respoder.SendOK(w, alphabeticallyHeadersKeys, alphabeticallyHeadersValues)
		return
	}
}

// sortHeadersCharactersAlphaBetically get headers from request and arrage every headers key and value in alphabetical order
// for example, if a header values is json it will set it to jnos
func (s *Server) sortHeadersCharactersAlphaBetically() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		alphabeticallyHeadersKeys := []string{}
		alphabeticallyHeadersValues := []string{}
		for name, values := range r.Header {
			for _, value := range values {
				alphabeticallyHeadersKeys = append(alphabeticallyHeadersKeys, SortString(name))
				alphabeticallyHeadersValues = append(alphabeticallyHeadersValues, SortString(value))
			}
		}
		sort.Strings(alphabeticallyHeadersKeys)
		sort.Strings(alphabeticallyHeadersValues)
		s.Respoder.SendOK(w, alphabeticallyHeadersKeys, alphabeticallyHeadersValues)
		return
	}
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
