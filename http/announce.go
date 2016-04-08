package http

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/mrd0ll4r/poke"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// BaseAnnounceResponse contains the fields present in all announce responses.
type BaseAnnounceResponse struct {
	FailureReason  string `bencode:"failure reason"`
	WarningMessage string `bencode:"warning message"`
	Interval       int    `bencode:"interval"`
	MinInterval    int    `bencode:"min interval"`
	Complete       int    `bencode:"complete"`
	Incomplete     int    `bencode:"incomplete"`
}

// CompactAnnounceResponse is a template to parse a compact bencoded announce
// response into.
type CompactAnnounceResponse struct {
	BaseAnnounceResponse
	Peers  []byte `bencode:"peers"`
	Peers6 []byte `bencode:"peers6"`
}

// Peer is a template to parse bencoded peer information into.
type Peer struct {
	ID   string `bencode:"peer id"`
	IP   string `bencode:"ip"`
	Port uint16 `bencode:"port"`
}

// NonCompactAnnounceResponse is a template to parse a non-compact bencoded
// announce response into.
type NonCompactAnnounceResponse struct {
	BaseAnnounceResponse
	Peers []Peer `bencode:"peers"`
}

// Client is a client for an http tracker.
type Client struct {
	address *url.URL
	client  *http.Client
}

var _ poke.Announcer = &Client{}

// NewClient returns a new client for the given announce URI.
func NewClient(announceAddress string) (*Client, error) {
	u, err := url.Parse(announceAddress)
	if err != nil {
		return nil, err
	}

	return &Client{
		address: u,
		client:  &http.Client{},
	}, nil
}

// Announce announces to the tracker.
func (c *Client) Announce(a poke.AnnounceRequest) (*poke.AnnounceResponse, error) {
	u, err := url.Parse(c.address.String())
	if err != nil {
		panic("url re-parse error")
	}
	v := u.Query()
	v.Set("info_hash", url.QueryEscape(string(a.InfoHash)))
	v.Set("peer_id", a.ID)
	v.Set("port", fmt.Sprint(a.Port))
	v.Set("uploaded", fmt.Sprint(a.Uploaded))
	v.Set("downloaded", fmt.Sprint(a.Downloaded))
	v.Set("left", fmt.Sprint(a.Left))
	if a.Compact {
		v.Set("compact", "1")
	}

	switch a.Event {
	case poke.EventStarted:
		v.Set("event", "started")
	case poke.EventStopped:
		v.Set("event", "stopped")
	case poke.EventCompleted:
		v.Set("event", "completed")
	default:
		return nil, errors.New("invalid event")
	}

	for _, b := range []byte(a.IP) {
		if b != 0 {
			v.Set("ip", a.IP.String())
		}
	}

	if a.Numwant != 0 {
		v.Set("numwant", fmt.Sprint(a.Numwant))
	}

	u.RawQuery = v.Encode()
	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response:", string(b))

	if a.Compact {
		r := CompactAnnounceResponse{}
		err = bencode.DecodeBytes(b, &r)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if r.FailureReason != "" {
			log.Println(r.FailureReason)
			return nil, errors.New(r.FailureReason)
		}
		if r.WarningMessage != "" {
			log.Println(r.FailureReason)
			return nil, errors.New(r.WarningMessage)
		}

		ann := poke.AnnounceResponse{
			Interval:    r.Interval,
			MinInterval: r.MinInterval,
			Incomplete:  r.Incomplete,
			Complete:    r.Complete,
			Peers:       make([]poke.Peer, 0),
		}

		for i := 0; i < len(r.Peers); i += 6 {
			peer := poke.Peer{}
			peer.IP = net.IPv4(r.Peers[0+i], r.Peers[1+i], r.Peers[2+i], r.Peers[3+i])
			reader := bytes.NewBuffer(r.Peers[4+i : 6+i])
			err = binary.Read(reader, binary.BigEndian, &peer.Port)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			ann.Peers = append(ann.Peers, peer)
		}

		for i := 0; i < len(r.Peers6); i += 18 {
			peer := poke.Peer{}
			peer.IP = net.IP(r.Peers6[i+0 : i+16])
			reader := bytes.NewBuffer(r.Peers[i+16 : i+18])
			err = binary.Read(reader, binary.BigEndian, &peer.Port)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			ann.Peers = append(ann.Peers, peer)
		}

		return &ann, nil
	}

	r := NonCompactAnnounceResponse{}
	err = bencode.DecodeBytes(b, &r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if r.FailureReason != "" {
		log.Println(r.FailureReason)
		return nil, errors.New(r.FailureReason)
	}
	if r.WarningMessage != "" {
		log.Println(r.WarningMessage)
		return nil, errors.New(r.WarningMessage)
	}

	ann := poke.AnnounceResponse{
		Interval:    r.Interval,
		MinInterval: r.MinInterval,
		Incomplete:  r.Incomplete,
		Complete:    r.Complete,
		Peers:       make([]poke.Peer, 0),
	}

	for _, peer := range r.Peers {
		p := poke.Peer{
			ID:   peer.ID,
			Port: peer.Port,
			IP:   net.ParseIP(peer.IP),
		}
		ann.Peers = append(ann.Peers, p)
	}

	return &ann, nil
}
