package poke

import "net"

// Event represents an event provided in a BitTorrent announce.
type Event int

// Events for announces.
const (
	EventStarted = iota
	EventStopped
	EventCompleted
)

// InfoHash represents a 20-byte infohash in hexadecimal format.
type InfoHash string

// Peer represents a peer in a BitTorrent swarm.
type Peer struct {
	ID   string
	Port uint16
	IP   net.IP
}

// AnnounceRequest represents an announce request.
type AnnounceRequest struct {
	InfoHash   InfoHash
	Uploaded   int
	Downloaded int
	Left       int
	Compact    bool
	Event      Event
	Numwant    int
	Peer
}

// AnnounceResponse represents an announce response.
type AnnounceResponse struct {
	Interval    int
	MinInterval int
	Complete    int
	Incomplete  int
	Peers       []Peer
}

// ScrapeRequest respresents a scrape request.
type ScrapeRequest struct {
	InfoHashes []InfoHash
}

// Scrape represents the set of information scraped for a single infohash.
type Scrape struct {
	Complete   int
	Downloaded int
	Incomplete int
}

// ScrapeResponse represents a scrape response.
type ScrapeResponse struct {
	Files map[InfoHash]Scrape
}

// Announcer provides the Announce method.
type Announcer interface {
	Announce(AnnounceRequest) (*AnnounceResponse, error)
}

// Scraper provides the Scrape method.
type Scraper interface {
	Scrape(ScrapeRequest) (*ScrapeResponse, error)
}
