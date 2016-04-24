package tests

import (
	"errors"
	"math/rand"
	"time"

	"github.com/mrd0ll4r/poke"
	"github.com/mrd0ll4r/poke/http"
)

// BasicHTTPAnnounce performs a basic, compact HTTP announce.
func BasicHTTPAnnounce(announceURI string) error {
	return basicHTTPAnnounce(announceURI)
}

// BasicHTTPNonCompactAnnounce performs a basic, non-compact HTTP announce.
func BasicHTTPNonCompactAnnounce(announceURI string) error {
	return basicHTTPNonCompactAnnounce(announceURI)
}

// BasicHTTPSeederAnnounce performs a basic HTTP announce as a seeder.
func BasicHTTPSeederAnnounce(announceURI string) error {
	return basicHTTPSeederAnnounce(announceURI)
}

// TrackerSupportsIPSpoofingHTTPAnnounce reports whether the tracker supports
// IP spoofing via HTTP announce.
func TrackerSupportsIPSpoofingHTTPAnnounce(announceURI string) (bool, error) {
	return trackerSupportsIPSpoofingHTTPAnnounce(announceURI)
}

// TrackerSupportsCompactHTTPAnnounce reports whether the tracker supports
// compact HTTP announces.
func TrackerSupportsCompactHTTPAnnounce(announceURI string) (bool, error) {
	return trackerSupportsCompactHTTPAnnounce(announceURI)
}

// TrackerSupportsNonCompactHTTPAnnounce reports whether the tracker supports
// non-compact HTTP announces.
func TrackerSupportsNonCompactHTTPAnnounce(announceURI string) (bool, error) {
	return trackerSupportsNonCompactHTTPAnnounce(announceURI)
}

// TrackerSupportsOptimizedSeederHTTPAnnounce reports whether the tracker
// supports optimizing seeder's HTTP announces by not returning other seeders
// in the peer list.
func TrackerSupportsOptimizedSeederHTTPAnnounce(announceURI string) (bool, error) {
	return trackerSupportsOptimizedSeederHTTPAnnounce(announceURI)
}

// TrackerSupportsAnnouncingPeerNotInPeerListHTTPAnnounce reports whether the
// tracker supports leaving the announcing peer out of the peer list returned.
func TrackerSupportsAnnouncingPeerNotInPeerListHTTPAnnounce(announceURI string) (bool, error) {
	return trackerSupportsAnnouncingPeerNotInPeerListHTTPAnnounce(announceURI)
}

// CheckReturnedPeersHTTPAnnounce checks whether the peers returned by an
// HTTP announce are correct.
func CheckReturnedPeersHTTPAnnounce(announceURI string) error {
	return checkReturnedPeersHTTPAnnounce(announceURI)
}

// InvalidShortInfohashHTTPAnnounce checks whether the tracker rejects announces
// with an invalid too short infohash.
func InvalidShortInfohashHTTPAnnounce(announceURI string) error {
	return invalidShortInfohashHTTPAnnounce(announceURI)
}

// InvalidLongInfohashHTTPAnnounce checks whether the tracker rejects announces
// with an invalid too long infohash.
func InvalidLongInfohashHTTPAnnounce(announceURI string) error {
	return invalidLongInfohashHTTPAnnounce(announceURI)
}

// InvalidShortPeerIDHTTPAnnounce checks whether the tracker rejects announces
// with an invalid too short peer ID.
func InvalidShortPeerIDHTTPAnnounce(announceURI string) error {
	return invalidShortPeerIDHTTPAnnounce(announceURI)
}

// InvalidLongPeerIDHTTPAnnounce checks whether the tracker rejects announces
// with an invalid too short peer ID.
func InvalidLongPeerIDHTTPAnnounce(announceURI string) error {
	return invalidLongPeerIDHTTPAnnounce(announceURI)
}

// InvalidNegativeUploadedHTTPAnnounce checks whether the tracker rejects
// announces with an invalid negative amount of bytes uploaded.
func InvalidNegativeUploadedHTTPAnnounce(announceURI string) error {
	return invalidNegativeUploadedHTTPAnnounce(announceURI)
}

// InvalidNegativeDownloadedHTTPAnnounce checks whether the tracker rejects
// announces with an invalid negative amount of bytes downloaded.
func InvalidNegativeDownloadedHTTPAnnounce(announceURI string) error {
	return invalidNegativeDownloadedHTTPAnnounce(announceURI)
}

// InvalidNegativeLeftHTTPAnnounce checks whether the tracker rejects
// announces with an invalid negative amount of bytes left.
func InvalidNegativeLeftHTTPAnnounce(announceURI string) error {
	return invalidNegativeLeftHTTPAnnounce(announceURI)
}

// InvalidEventHTTPAnnounce checks whether the tracker rejects announces with
// an invalid event.
func InvalidEventHTTPAnnounce(announceURI string) error {
	return invalidEventHTTPAnnounce(announceURI)
}

func basicHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if resp.Complete != 0 || resp.Incomplete != 1 {
		return errors.New("first announce does not have zero seeders and one leecher (the one that just announced)")
	}

	if len(resp.Peers) > 0 {
		return errors.New("first announce is not empty")
	}

	return nil
}

func basicHTTPNonCompactAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  false,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if resp.Complete != 0 || resp.Incomplete != 1 {
		return errors.New("first announce does not have zero seeders and one leecher (the one that just announced)")
	}

	if len(resp.Peers) > 0 {
		return errors.New("first announce is not empty")
	}

	return nil
}

func trackerSupportsCompactHTTPAnnounce(announceURI string) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, poke.WrapError("unable to create client", err)
	}

	leecher1 := poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano())))
	leecher2 := poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano())))

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	req.Peer = leecher2

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 1 {
		return false, errors.New("announce did not return the other known leecher")
	}

	if resp.Peers[0].ID == "" {
		return true, nil
	}

	return false, nil
}

func trackerSupportsNonCompactHTTPAnnounce(announceURI string) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, poke.WrapError("unable to create client", err)
	}

	leecher1 := poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano())))
	leecher2 := poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano())))

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  false,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	req.Peer = leecher2

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 1 {
		return false, errors.New("announce did not return the other known leecher")
	}

	if resp.Peers[0].ID == leecher1.ID {
		return true, nil
	}

	return false, nil
}

func basicHTTPSeederAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     0,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if resp.Complete != 1 || resp.Incomplete != 0 {
		return errors.New("first announce does not have zero leechers and one seeder (the one that just announced)")
	}

	if len(resp.Peers) > 0 {
		return errors.New("first announce is not empty")
	}

	return nil
}

func trackerSupportsAnnouncingPeerNotInPeerListHTTPAnnounce(announceURI string) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	for _, p := range resp.Peers {
		if p.IsEqual(req.Peer) {
			return false, nil
		}
	}

	req.Left = 50
	req.Event = poke.EventNone

	resp, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) == 0 {
		return true, nil
	} else if len(resp.Peers) == 1 {
		if resp.Peers[0].IsEqual(req.Peer) {
			return false, nil
		}
	}
	return false, errors.New("second announce with equal peer did return more than one peer")

}

func trackerSupportsIPSpoofingHTTPAnnounce(announceURI string) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, poke.WrapError("unable to create client", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	leecher1 := poke.NewPeer(r)
	leecher2 := poke.NewPeer(r)

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(r),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  false,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	req.Peer = leecher2
	req.Left = 120

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 1 || resp.Peers[0].ID != leecher1.ID {
		return false, errors.New("announce of second peer did not return known peer")
	}

	if resp.Peers[0].IP.Equal(leecher1.IP) {
		return true, nil
	}
	return false, nil
}

func trackerSupportsOptimizedSeederHTTPAnnounce(announceURI string) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, poke.WrapError("unable to create client", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	leecher1 := poke.NewPeer(r)
	seeder1 := poke.NewPeer(r)
	seeder2 := poke.NewPeer(r)

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(r),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  false,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	req.Peer = seeder1
	req.Left = 0

	_, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	req.Peer = seeder2

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) < 1 {
		return false, errors.New("announce did not return an expected peer")
	} else if len(resp.Peers) > 1 {
		return false, nil
	}

	if resp.Peers[0].ID == leecher1.ID {
		return true, nil
	}

	return false, nil
}

func checkReturnedPeersHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	leecher1 := poke.NewPeer(r)
	leecher2 := poke.NewPeer(r)
	seeder1 := poke.NewPeer(r)
	seeder2 := poke.NewPeer(r)

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(r),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  false,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	req.Peer = leecher2
	req.Left = 120

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 1 || !resp.Peers[0].IsEqual(leecher1) {
		return errors.New("announce did not return the other known peer")
	}

	req.Peer = seeder1
	req.Left = 0

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 2 {
		return errors.New("seeder announce did not return the two known peers")
	}

	if !((resp.Peers[0].IsEqual(leecher1) || resp.Peers[0].IsEqual(leecher2)) && (resp.Peers[1].IsEqual(leecher1) || resp.Peers[1].IsEqual(leecher2))) {
		return errors.New("seeder announce did not return the two known peers")
	}

	req.Peer = seeder2

	_, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	req.Peer = leecher1
	req.Left = 80
	req.Event = poke.EventNone

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	if len(resp.Peers) != 3 {
		return errors.New("leecher announce did not return both known seeders and the other leecher")
	}

	if !((resp.Peers[0].IsEqual(leecher2) || resp.Peers[0].IsEqual(seeder1) || resp.Peers[0].IsEqual(seeder2)) &&
		(resp.Peers[1].IsEqual(leecher2) || resp.Peers[1].IsEqual(seeder1) || resp.Peers[1].IsEqual(seeder2)) &&
		(resp.Peers[2].IsEqual(leecher2) || resp.Peers[2].IsEqual(seeder1) || resp.Peers[2].IsEqual(seeder2))) {
		return errors.New("leecher announce did not return both known seeders and the other leecher")
	}

	return nil
}

func invalidShortInfohashHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.InfoHash([]byte{30, 30, 30}),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid too short infohash did not fail")
	}

	return nil
}

func invalidLongInfohashHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.InfoHash([]byte{
			30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
			30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
			30}),
		Peer:    poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:   poke.EventStarted,
		Numwant: 50,
		Compact: true,
		Left:    100,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid too long infohash did not fail")
	}

	return nil
}

func invalidShortPeerIDHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	req.Peer.ID = req.Peer.ID[:19]

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid too short peerID did not fail")
	}

	return nil
}

func invalidLongPeerIDHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	req.Peer.ID += "9"

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid too long peerID did not fail")
	}

	return nil
}

func invalidNegativeUploadedHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     100,
		Uploaded: -1,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid negative uploaded did not fail")
	}

	return nil
}

func invalidNegativeDownloadedHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash:   poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:       poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:      poke.EventStarted,
		Numwant:    50,
		Compact:    true,
		Left:       100,
		Downloaded: -1,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid negative downloaded did not fail")
	}

	return nil
}

func invalidNegativeLeftHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Compact:  true,
		Left:     -1,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid negative left did not fail")
	}

	return nil
}

func invalidEventHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return poke.WrapError("unable to create client", err)
	}

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventInvalid,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	_, err = c.Announce(req)
	if err == nil {
		return errors.New("announce with invalid event did not fail")
	}

	return nil
}
