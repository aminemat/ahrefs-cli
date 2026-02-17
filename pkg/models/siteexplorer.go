package models

// DomainRatingResponse represents the domain rating API response
type DomainRatingResponse struct {
	DomainRating DomainRating `json:"domain_rating"`
}

// DomainRating contains the domain rating value
type DomainRating struct {
	DomainRating float64 `json:"domain_rating"`
}

// BacklinksStatsResponse represents the backlinks stats API response
type BacklinksStatsResponse struct {
	Metrics BacklinksMetrics `json:"metrics"`
}

// BacklinksMetrics contains backlink metrics
type BacklinksMetrics struct {
	Live         int `json:"live"`
	Refdomains   int `json:"refdomains,omitempty"`
	DoFollow     int `json:"dofollow,omitempty"`
	Governmental int `json:"governmental,omitempty"`
	Educational  int `json:"educational,omitempty"`
}

// BacklinksResponse represents a list of backlinks
type BacklinksResponse struct {
	Backlinks []Backlink `json:"backlinks"`
}

// Backlink represents a single backlink
type Backlink struct {
	URLFrom      string `json:"url_from"`
	URLTo        string `json:"url_to"`
	DomainRating int    `json:"domain_rating,omitempty"`
	AhrefsRank   int    `json:"ahrefs_rank,omitempty"`
	Anchor       string `json:"anchor,omitempty"`
	HTTPCode     int    `json:"http_code,omitempty"`
	FirstSeen    string `json:"first_seen,omitempty"`
	LastVisited  string `json:"last_visited,omitempty"`
}

// RefDomainsResponse represents a list of referring domains
type RefDomainsResponse struct {
	RefDomains []RefDomain `json:"refdomains"`
}

// RefDomain represents a single referring domain
type RefDomain struct {
	Domain       string `json:"domain"`
	DomainRating int    `json:"domain_rating,omitempty"`
	Backlinks    int    `json:"backlinks,omitempty"`
	FirstSeen    string `json:"first_seen,omitempty"`
}
