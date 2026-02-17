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
	URLFrom      string  `json:"url_from"`
	URLTo        string  `json:"url_to"`
	DomainRating float64 `json:"domain_rating,omitempty"`
	AhrefsRank   int     `json:"ahrefs_rank,omitempty"`
	Anchor       string  `json:"anchor,omitempty"`
	HTTPCode     int     `json:"http_code,omitempty"`
	FirstSeen    string  `json:"first_seen,omitempty"`
	LastVisited  string  `json:"last_visited,omitempty"`
	LinkType     string  `json:"link_type,omitempty"`
	URLRating    float64 `json:"url_rating,omitempty"`
	Traffic      int     `json:"traffic,omitempty"`
}

// RefDomainsResponse represents a list of referring domains
type RefDomainsResponse struct {
	RefDomains []RefDomain `json:"refdomains"`
}

// RefDomain represents a single referring domain
type RefDomain struct {
	Domain       string  `json:"domain"`
	DomainRating float64 `json:"domain_rating,omitempty"`
	URLRating    float64 `json:"url_rating,omitempty"`
	AhrefsRank   int     `json:"ahrefs_rank,omitempty"`
	Backlinks    int     `json:"backlinks,omitempty"`
	DoFollow     int     `json:"dofollow,omitempty"`
	LinkedPages  int     `json:"linked_pages,omitempty"`
	FirstSeen    string  `json:"first_seen,omitempty"`
	LastVisited  string  `json:"last_visited,omitempty"`
}

// AnchorsResponse represents a list of anchor texts
type AnchorsResponse struct {
	Anchors []Anchor `json:"anchors"`
}

// Anchor represents a single anchor text entry
type Anchor struct {
	Anchor      string `json:"anchor"`
	Backlinks   int    `json:"backlinks,omitempty"`
	Refdomains  int    `json:"refdomains,omitempty"`
	FirstSeen   string `json:"first_seen,omitempty"`
	LastVisited string `json:"last_visited,omitempty"`
}

// OrganicKeywordsResponse represents a list of organic keywords
type OrganicKeywordsResponse struct {
	Keywords []OrganicKeyword `json:"keywords"`
}

// OrganicKeyword represents a single organic keyword entry
type OrganicKeyword struct {
	Keyword      string  `json:"keyword"`
	Position     int     `json:"position,omitempty"`
	SearchVolume int     `json:"volume,omitempty"`
	Traffic      int     `json:"traffic,omitempty"`
	KD           float64 `json:"kd,omitempty"`
	URL          string  `json:"url,omitempty"`
	Country      string  `json:"country,omitempty"`
}

// TopPagesResponse represents a list of top pages
type TopPagesResponse struct {
	Pages []TopPage `json:"pages"`
}

// TopPage represents a single top page entry
type TopPage struct {
	URL          string  `json:"url"`
	Traffic      int     `json:"traffic,omitempty"`
	TrafficValue int     `json:"traffic_value,omitempty"`
	Keywords     int     `json:"keywords,omitempty"`
	TopKeyword   string  `json:"top_keyword,omitempty"`
	Position     int     `json:"position,omitempty"`
	Volume       int     `json:"volume,omitempty"`
	URLRating    float64 `json:"url_rating,omitempty"`
}

// BrokenBacklinksResponse represents a list of broken backlinks
type BrokenBacklinksResponse struct {
	Backlinks []BrokenBacklink `json:"backlinks"`
}

// BrokenBacklink represents a single broken backlink
type BrokenBacklink struct {
	URLFrom      string  `json:"url_from"`
	URLTo        string  `json:"url_to"`
	DomainRating float64 `json:"domain_rating,omitempty"`
	HTTPCode     int     `json:"http_code,omitempty"`
	Anchor       string  `json:"anchor,omitempty"`
	FirstSeen    string  `json:"first_seen,omitempty"`
	LastVisited  string  `json:"last_visited,omitempty"`
}

// LinkedDomainsResponse represents a list of linked domains
type LinkedDomainsResponse struct {
	LinkedDomains []LinkedDomain `json:"linked_domains"`
}

// LinkedDomain represents a single linked domain
type LinkedDomain struct {
	Domain       string  `json:"domain"`
	DomainRating float64 `json:"domain_rating,omitempty"`
	LinkedPages  int     `json:"linked_pages,omitempty"`
	Backlinks    int     `json:"backlinks,omitempty"`
	FirstSeen    string  `json:"first_seen,omitempty"`
}

// MetricsResponse represents site metrics
type MetricsResponse struct {
	Metrics SiteMetrics `json:"metrics"`
}

// SiteMetrics contains comprehensive site metrics
type SiteMetrics struct {
	OrgKeywords      int     `json:"org_keywords,omitempty"`
	OrgKeywords2     int     `json:"org_keywords_2,omitempty"`
	OrgTraffic       int     `json:"org_traffic,omitempty"`
	OrgCost          float64 `json:"org_cost,omitempty"`
	PaidKeywords     int     `json:"paid_keywords,omitempty"`
	PaidTraffic      int     `json:"paid_traffic,omitempty"`
	PaidCost         float64 `json:"paid_cost,omitempty"`
	FeaturedSnippets int     `json:"featured_snippets,omitempty"`
}

// MetricsHistoryResponse represents historical metrics data
type MetricsHistoryResponse struct {
	Metrics []MetricsHistoryEntry `json:"metrics"`
}

// MetricsHistoryEntry represents a single historical metrics entry
type MetricsHistoryEntry struct {
	Date         string  `json:"date"`
	OrgKeywords  int     `json:"org_keywords,omitempty"`
	OrgTraffic   int     `json:"org_traffic,omitempty"`
	OrgCost      float64 `json:"org_cost,omitempty"`
	PaidKeywords int     `json:"paid_keywords,omitempty"`
	PaidTraffic  int     `json:"paid_traffic,omitempty"`
	DomainRating float64 `json:"domain_rating,omitempty"`
}

// PagesByTrafficResponse represents pages sorted by traffic
type PagesByTrafficResponse struct {
	Pages []PageByTraffic `json:"pages"`
}

// PageByTraffic represents a page with traffic data
type PageByTraffic struct {
	URL          string  `json:"url"`
	Traffic      int     `json:"traffic,omitempty"`
	TrafficValue int     `json:"traffic_value,omitempty"`
	Keywords     int     `json:"keywords,omitempty"`
	URLRating    float64 `json:"url_rating,omitempty"`
}

// BestByLinksResponse represents pages sorted by backlinks
type BestByLinksResponse struct {
	Pages []PageByLinks `json:"pages"`
}

// PageByLinks represents a page with link data
type PageByLinks struct {
	URL        string  `json:"url"`
	Backlinks  int     `json:"backlinks,omitempty"`
	Refdomains int     `json:"refdomains,omitempty"`
	URLRating  float64 `json:"url_rating,omitempty"`
	Traffic    int     `json:"traffic,omitempty"`
	FirstSeen  string  `json:"first_seen,omitempty"`
}
