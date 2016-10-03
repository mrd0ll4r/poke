package tests

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/mrd0ll4r/poke"
	"github.com/mrd0ll4r/poke/http"
	"github.com/mrd0ll4r/poke/udp"
)

// Test represents a test performed on a tracker.
type Test struct {
	Name         string
	Run          bool
	NotRunReason string
	Result       TestResult
}

// TestResult represents the result of a test.
type TestResult struct {
	Result interface{}
	Err    error
}

// HTTPResult represents the result of all tests performed on an HTTP tracker.
type HTTPResult struct {
	TrackerResult
	SupportsCompact    bool
	SupportsNonCompact bool
}

// TrackerResult represents the result of all tests performed on a tracker.
type TrackerResult struct {
	SupportsAnnouncingPeerNotInPeerList bool
	SupportsIPSpoofing                  bool
	SupportsOptimizedSeederResponse     bool
	Tests                               []Test
}

// UDPResult represents the result of all tests performed on a UDP tracker.
type UDPResult struct {
	TrackerResult
}

// TestUDPTracker runs tests on a UDP tracker to determine its functionality
// and feature-completeness.
func TestUDPTracker(addr string) (*UDPResult, error) {
	toReturn := &UDPResult{
		TrackerResult: TrackerResult{
			Tests: make([]Test, 0),
		},
	}

	f := func() (poke.Announcer, error) {
		c, err := udp.NewClient(addr)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	err := runAll(f, &toReturn.TrackerResult)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func testTrackerSupportsAnnouncingPeerNotInPeerList(c poke.Announcer, result *TrackerResult) error {
	t := Test{
		Name: "trackerSupportsAnnouncingPeerNotInPeerListAnnounce",
	}

	res, err := trackerSupportsAnnouncingPeerNotInPeerListAnnounce(c)
	t.Run = true
	t.Result.Err = err
	t.Result.Result = res
	if err != nil {
		result.SupportsAnnouncingPeerNotInPeerList = false
	} else {
		result.SupportsAnnouncingPeerNotInPeerList = res
	}
	result.Tests = append(result.Tests, t)

	return nil
}

func testTrackerSupportsIPSpoofing(c poke.Announcer, result *TrackerResult) error {
	t := Test{
		Name: "trackerSupportsIPSpoofingAnnounce",
	}

	res, err := trackerSupportsIPSpoofingAnnounce(c, result.SupportsAnnouncingPeerNotInPeerList)
	t.Run = true
	t.Result.Err = err
	t.Result.Result = res
	if err != nil {
		result.SupportsIPSpoofing = false
	} else {
		result.SupportsIPSpoofing = res
	}
	result.Tests = append(result.Tests, t)

	return nil
}

func testTrackerSupportsOptimizedSeederResponse(c poke.Announcer, result *TrackerResult) error {
	t := Test{
		Name: "trackerSupportsOptimizedSeederResponseAnnounce",
	}
	res, err := trackerSupportsOptimizedSeederAnnounce(c)
	t.Run = true
	t.Result.Err = err
	t.Result.Result = res
	if err != nil {
		result.SupportsOptimizedSeederResponse = false
	} else {
		result.SupportsOptimizedSeederResponse = res
	}
	result.Tests = append(result.Tests, t)

	return nil
}

func runBasicAnnounce(c poke.Announcer, result *TrackerResult) error {
	t := Test{
		Name: "basicAnnounce",
	}
	err := basicAnnounce(c, result.SupportsAnnouncingPeerNotInPeerList)
	t.Run = true
	t.Result.Err = err
	result.Tests = append(result.Tests, t)

	return nil
}

func runAll(f func() (poke.Announcer, error), result *TrackerResult) error {
	c, err := f()
	if err != nil {
		return err
	}
	err = testTrackerSupportsAnnouncingPeerNotInPeerList(c, result)
	if err != nil {
		return err
	}

	err = testTrackerSupportsIPSpoofing(c, result)
	if err != nil {
		return err
	}

	err = testTrackerSupportsOptimizedSeederResponse(c, result)
	if err != nil {
		return err
	}

	err = runBasicAnnounce(c, result)

	return err
}

// TestHTTPTracker runs tests on an HTTP tracker to determine its functionality
// and feature-completeness.
func TestHTTPTracker(announceURI string) (*HTTPResult, error) {
	toReturn := &HTTPResult{
		TrackerResult: TrackerResult{
			Tests: make([]Test, 0),
		},
	}
	t := Test{
		Name: "trackerSupportsCompactAnnounce",
		Run:  true,
	}
	supportsCompact, err := trackerSupportsCompactHTTPAnnounce(announceURI)
	t.Result.Err = err
	t.Result.Result = supportsCompact
	if err != nil {
		toReturn.SupportsCompact = false
	} else {
		toReturn.SupportsCompact = supportsCompact
	}
	toReturn.Tests = append(toReturn.Tests, t)

	t = Test{
		Name: "trackerSupportsNonCompactAnnounce",
		Run:  true,
	}
	supportsNonCompact, err := trackerSupportsNonCompactHTTPAnnounce(announceURI)
	t.Result.Err = err
	t.Result.Result = supportsNonCompact
	if err != nil {
		toReturn.SupportsNonCompact = false
	} else {
		toReturn.SupportsNonCompact = supportsNonCompact
	}
	toReturn.Tests = append(toReturn.Tests, t)

	if !toReturn.SupportsCompact && !toReturn.SupportsNonCompact {
		// We cannot run tests.
		return toReturn, nil
	}

	f := func() (poke.Announcer, error) {
		c, err := http.NewClient(announceURI)
		if err != nil {
			return nil, err
		}
		c.OverrideCompact(toReturn.SupportsCompact)
		return c, nil
	}

	err = runAll(f, &toReturn.TrackerResult)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

// BasicHTTPCompactAnnounce performs a basic, compact HTTP announce.
func BasicHTTPCompactAnnounce(announceURI string, trackerSupportsAnnouncingPeerNotInPeerList bool) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return err
	}
	c.OverrideCompact(true)
	return basicAnnounce(c, trackerSupportsAnnouncingPeerNotInPeerList)
}

// BasicHTTPNonCompactAnnounce performs a basic, non-compact HTTP announce.
func BasicHTTPNonCompactAnnounce(announceURI string, trackerSupportsAnnouncingPeerNotInPeerList bool) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return err
	}
	c.OverrideCompact(false)
	return basicAnnounce(c, trackerSupportsAnnouncingPeerNotInPeerList)
}

// BasicHTTPSeederAnnounce performs a basic HTTP announce as a seeder.
func BasicHTTPSeederAnnounce(announceURI string) error {
	return basicHTTPSeederAnnounce(announceURI)
}

// TrackerSupportsIPSpoofingHTTPNonCompactAnnounce reports whether the tracker
// supports IP spoofing via HTTP non-compact announce.
func TrackerSupportsIPSpoofingHTTPNonCompactAnnounce(announceURI string, trackerSupportsOptimizedPeerList bool) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, err
	}
	c.OverrideCompact(false)
	return trackerSupportsIPSpoofingAnnounce(c, trackerSupportsOptimizedPeerList)
}

// TrackerSupportsIPSpoofingHTTPCompactAnnounce reports whether the tracker
// supports IP spoofing via HTTP compact announce.
func TrackerSupportsIPSpoofingHTTPCompactAnnounce(announceURI string, trackerSupportsOptimizedPeerList bool) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, err
	}
	c.OverrideCompact(true)
	return trackerSupportsIPSpoofingAnnounce(c, trackerSupportsOptimizedPeerList)
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
// in the peer list for compact announces.
func TrackerSupportsOptimizedSeederHTTPAnnounce(announceURI string, trackerSupportsCompactAnnounce bool) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, err
	}
	c.OverrideCompact(trackerSupportsCompactAnnounce)
	return trackerSupportsOptimizedSeederAnnounce(c)
}

// TrackerSupportsAnnouncingPeerNotInPeerListHTTPCompactAnnounce reports
// whether the tracker supports leaving the announcing peer out of the peer list
// returned for compact announces.
func TrackerSupportsAnnouncingPeerNotInPeerListHTTPCompactAnnounce(announceURI string, trackerSupportsCompactAnnounce bool) (bool, error) {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return false, err
	}
	c.OverrideCompact(trackerSupportsCompactAnnounce)
	return trackerSupportsAnnouncingPeerNotInPeerListAnnounce(c)
}

// CheckReturnedPeersHTTPAnnounce checks whether the peers returned by an
// HTTP announce are correct.
func CheckReturnedPeersHTTPAnnounce(announceURI string) error {
	c, err := http.NewClient(announceURI)
	if err != nil {
		return err
	}
	return checkReturnedPeersAnnounce(c)
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
	c, err := http.NewClient(announceURI)
	if err != nil {
		return err
	}
	return invalidEventAnnounce(c)
}

func basicAnnounce(c poke.Announcer, trackerSupportsAnnouncingPeerNotInPeerList bool) error {
	if poke.Debug {
		log.Println("Running basicAnnounce")
	}
	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if resp.Complete != 0 || resp.Incomplete != 1 {
			return errors.New("first announce does not have zero seeders and one leecher (the one that just announced)")
		}

		if trackerSupportsAnnouncingPeerNotInPeerList {
			if len(resp.Peers) != 0 {
				return errors.New("first announce is not empty")
			}
		} else {
			if len(resp.Peers) != 1 {
				return errors.New("expected one peer for first announce")
			}
		}
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
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

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		switch len(resp.Peers) {
		case 2:
			if resp.Peers[0].Port == leecher1.Port || resp.Peers[1].Port == leecher1.Port {
				return true, nil
			}
			return false, errors.New("announce returned unknown peer")
		case 1:
			if resp.Peers[0].Port == leecher1.Port {
				return true, nil
			}
			return false, errors.New("announce returned unknown peer")
		case 0:
			return false, errors.New("announce did not return the other known leecher")
		default:
			return false, errors.New("announce returned too many peers")
		}
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
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

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		switch len(resp.Peers) {
		case 2:
			if resp.Peers[0].Port == leecher1.Port || resp.Peers[1].Port == leecher1.Port {
				return true, nil
			}
			return false, errors.New("announce returned unknown peer")
		case 1:
			if resp.Peers[0].Port == leecher1.Port {
				return true, nil
			}
			return false, errors.New("announce returned unknown peer")
		case 0:
			return false, errors.New("announce did not return the other known leecher")
		default:
			return false, errors.New("announce returned too many peers")
		}
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
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
	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if resp.Complete != 1 || resp.Incomplete != 0 {
			return errors.New("first announce does not have zero leechers and one seeder (the one that just announced)")
		}

		if len(resp.Peers) > 0 {
			return errors.New("first announce is not empty")
		}
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
	}

	return nil
}

func trackerSupportsAnnouncingPeerNotInPeerListAnnounce(c poke.Announcer) (bool, error) {
	if poke.Debug {
		log.Println("Running trackerSupportsAnnouncingPeerNotInPeerList")
	}
	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventStarted,
		Numwant:  50,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) > 0 {
			return false, nil
		}
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	}

	req.Left = 50
	req.Event = poke.EventNone

	resp, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) == 0 {
			return true, nil
		} else if len(resp.Peers) == 1 {
			if resp.Peers[0].IsEqual(req.Peer) {
				return false, nil
			}
			return false, errors.New("second announce with equal peer returned unknown peer")
		}
		return false, errors.New("second announce with equal peer did return more than one peer")
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	}
	return false, nil
}

func trackerSupportsIPSpoofingAnnounce(c poke.Announcer, trackerSupportsAnnouncingPeerNotInPeerList bool) (bool, error) {
	if poke.Debug {
		log.Println("Running trackerSupportsIPSpoofingAnnounce")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	leecher1 := poke.NewPeer(r)
	leecher2 := poke.NewPeer(r)

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(r),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	default:
	}

	req.Peer = leecher2
	req.Left = 120

	resp, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if trackerSupportsAnnouncingPeerNotInPeerList {
			if len(resp.Peers) != 1 || resp.Peers[0].Port != leecher1.Port {
				return false, errors.New("announce of second peer did not return known peer")
			}

			if resp.Peers[0].IP.Equal(leecher1.IP) {
				return true, nil
			}
		} else {
			if len(resp.Peers) != 2 || (resp.Peers[0].Port != leecher1.Port && resp.Peers[0].Port != leecher2.Port) || (resp.Peers[1].Port != leecher1.Port && resp.Peers[1].Port != leecher2.Port) {
				return false, errors.New("announce of second peer returned unknown peer")
			}

			if resp.Peers[0].IP.Equal(leecher1.IP) || resp.Peers[1].IP.Equal(leecher2.IP) {
				return true, nil
			}
		}

	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	}
	return false, nil
}

func trackerSupportsOptimizedSeederAnnounce(c poke.Announcer) (bool, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	leecher1 := poke.NewPeer(r)
	seeder1 := poke.NewPeer(r)
	seeder2 := poke.NewPeer(r)

	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(r),
		Peer:     leecher1,
		Event:    poke.EventStarted,
		Numwant:  50,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	default:
	}

	req.Peer = seeder1
	req.Left = 0

	resp, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	default:
	}

	req.Peer = seeder2

	resp, err = c.Announce(req)
	if err != nil {
		return false, poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) < 1 {
			return false, errors.New("announce did not return an expected peer")
		} else if len(resp.Peers) > 1 {
			return false, nil
		}

		if resp.Peers[0].Port == leecher1.Port {
			return true, nil
		}
	case poke.ErrorResponse:
		return false, errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return false, errors.New("tracker returned warning: " + string(resp))
	}

	return false, nil
}

func checkReturnedPeersAnnounce(c poke.Announcer) error {
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
	default:
	}

	req.Peer = leecher2
	req.Left = 120

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) != 1 || !resp.Peers[0].IsEqual(leecher1) {
			return errors.New("announce did not return the other known peer")
		}
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
	}

	req.Peer = seeder1
	req.Left = 0

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}

	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) != 2 {
			return errors.New("seeder announce did not return the two known peers")
		}

		if !((resp.Peers[0].IsEqual(leecher1) || resp.Peers[0].IsEqual(leecher2)) && (resp.Peers[1].IsEqual(leecher1) || resp.Peers[1].IsEqual(leecher2))) {
			return errors.New("seeder announce did not return the two known peers")
		}
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
	}

	req.Peer = seeder2

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
	default:
	}

	req.Peer = leecher1
	req.Left = 80
	req.Event = poke.EventNone

	resp, err = c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp := resp.(type) {
	case poke.AnnounceResponse:
		if len(resp.Peers) != 3 {
			return errors.New("leecher announce did not return both known seeders and the other leecher")
		}

		if !((resp.Peers[0].IsEqual(leecher2) || resp.Peers[0].IsEqual(seeder1) || resp.Peers[0].IsEqual(seeder2)) &&
			(resp.Peers[1].IsEqual(leecher2) || resp.Peers[1].IsEqual(seeder1) || resp.Peers[1].IsEqual(seeder2)) &&
			(resp.Peers[2].IsEqual(leecher2) || resp.Peers[2].IsEqual(seeder1) || resp.Peers[2].IsEqual(seeder2))) {
			return errors.New("leecher announce did not return both known seeders and the other leecher")
		}
	case poke.ErrorResponse:
		return errors.New("tracker returned errror: " + string(resp))
	case poke.WarningResponse:
		return errors.New("tracker returned warning: " + string(resp))
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with too short infohash did not return error")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid too long infohash did not fail")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid too short peerID did not fail")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid too long peerID did not fail")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid negative uploaded did not fail")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid negative downloaded did not fail")
	}
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

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid negative left did not fail")
	}
}

func invalidEventAnnounce(c poke.Announcer) error {
	req := poke.AnnounceRequest{
		InfoHash: poke.NewInfohash(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Peer:     poke.NewPeer(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Event:    poke.EventInvalid,
		Numwant:  50,
		Compact:  true,
		Left:     100,
	}

	resp, err := c.Announce(req)
	if err != nil {
		return poke.WrapError("unable to perform announce", err)
	}
	switch resp.(type) {
	case poke.ErrorResponse:
		return nil
	default:
		return errors.New("announce with invalid event did not fail")
	}
}
