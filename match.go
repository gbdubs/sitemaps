package sitemaps

import (
	"math"
	"strings"
)

func (s *Sitemap) bestFuzzyMatch(target string) (string, int) {
	tl := simplifyStringForMatching(target)
	bestLev := math.MaxInt
	best := ""
	for ul, _ := range s.LastUpdated {
		sul := simplifyStringForMatching(ul)
		lev := levenshtein([]rune(tl), []rune(sul))
		if lev < bestLev {
			bestLev = lev
			best = ul
		}
	}
	return best, bestLev
}

func simplifyStringForMatching(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "www.", "")
	s = strings.ReplaceAll(s, "http://", "")
	s = strings.ReplaceAll(s, "https://", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}
