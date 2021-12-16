package sitemaps

import (
	"fmt"
	"strings"
	"time"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
)

// This is presently irrelevant since we use amass.Get()
// rather than amass.GetAll().
const sitemapMaxConcurrentRequests = 1

func getSitemapFromURL(url string) (*Sitemap, error) {
	sitemap := &Sitemap{
		URL: url,
	}
	req := amass.GetRequest{
		Site:                      "sitemap",
		RequestKey:                sitemapMemoizationKey(url),
		URL:                       url,
		SiteMaxConcurrentRequests: sitemapMaxConcurrentRequests,
		Attribution: attributions.Attribution{
			Context: []string{"Reading sitemap from " + url},
		},
	}
	resp, err := req.Get()
	if err != nil {
		return sitemap, fmt.Errorf("Couldn't retrieve sitemap at %s: %v", url, err)
	}
	sitemap.Attribution = resp.Attribution
	u := &urlSet{}
	err = resp.AsXMLObject(u)
	if err != nil {
		return sitemap, fmt.Errorf("Couldn't parse sitemap at %s: %v", url, err)
	}
	sitemap.LastUpdated, err = u.lastUpdatedMap()
	if err != nil {
		return sitemap, err
	}
	return sitemap, nil
}

func getPagedSitemapFromURL(url string) (*Sitemap, error) {
	sitemap := &Sitemap{
		URL:         url,
		LastUpdated: make(map[string]time.Time),
	}
	page := 1
	for {
		s, e := getSitemapFromURL(fmt.Sprintf("%s?page=%d", url, page))
		page++
		if e != nil {
			return sitemap, e
		}
		if len(s.LastUpdated) == 0 {
			return sitemap, nil
		}
		sitemap.Attribution = s.Attribution
		for k, v := range s.LastUpdated {
			sitemap.LastUpdated[k] = v
		}
	}
}

func sitemapMemoizationKey(url string) string {
	key := strings.ReplaceAll(url, "/", " ")
	key = strings.ReplaceAll(key, "www.", "")
	key = strings.ReplaceAll(key, "http:", "")
	key = strings.ReplaceAll(key, "https:", "")
	key = strings.ReplaceAll(key, "sitemap.xml", "")
	key = strings.ReplaceAll(key, ".", " ")
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "  ", "_")
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

type urlSet struct {
	URLs []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Location        string `xml:"loc"`
	ChangeFrequency string `xml:"changefreq"`
	Priority        string `xml:"priority"`
	LastModified    string `xml:"lastmod"`
}

func (s *urlSet) lastUpdatedMap() (map[string]time.Time, error) {
	m := make(map[string]time.Time)
	for _, u := range s.URLs {
		var t time.Time
		if u.LastModified != "" {
			tt, err := parseTime(u.LastModified)
			if err != nil {
				return m, err
			}
			t = tt
		}
		m[u.Location] = t
	}
	return m, nil
}

// Add more time formats here as we go along.
var timeFormats = [...]string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04Z"}

func parseTime(s string) (time.Time, error) {
	var err error
	for _, f := range timeFormats {
		t, e := time.Parse(f, s)
		if e == nil {
			return t, nil
		}
		err = e
	}
	return time.Time{}, err
}
