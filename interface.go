package sitemaps

import (
	"time"

	"github.com/gbdubs/attributions"
)

type Sitemap struct {
	URL         string
	Attribution attributions.Attribution
	LastUpdated map[string]time.Time
}

func GetSitemapFromURL(url string) (*Sitemap, error) {
	return getSitemapFromURL(url)
}

func (s *Sitemap) BestFuzzyMatch(target string) (string, int) {
	return s.bestFuzzyMatch(target)
}
