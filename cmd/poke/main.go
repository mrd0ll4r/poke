package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mrd0ll4r/poke"
	"github.com/mrd0ll4r/poke/tests"
)

func init() {
	flag.StringVar(&announceURI, "a", "http://tracker.org:6881/announce", "the announce URI")
	flag.StringVar(&udpAnnounceURI, "u", "tracker.org:6881", "the UDP announce URI")
	flag.BoolVar(&debug, "debug", false, "debug mode")
}

var (
	announceURI    string
	udpAnnounceURI string
	debug          bool
)

func main() {
	flag.Parse()

	poke.Debug = debug

	if f := flag.Lookup("u"); f != nil && f.Value.String() != f.DefValue {
		runUDPTests(udpAnnounceURI)
	} else {
		runHTTPTests(announceURI)
	}
}

func runUDPTests(addr string) {
	res, err := tests.TestUDPTracker(addr)
	if err != nil {
		log.Fatal(err)
	}

	formatTrackerResult(res.TrackerResult)
}

func runHTTPTests(announceURI string) {
	res, err := tests.TestHTTPTracker(announceURI)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tracker supports HTTP compact announces: %t\n", res.SupportsCompact)
	fmt.Printf("Tracker supports HTTP non-compact announces: %t\n", res.SupportsNonCompact)
	formatTrackerResult(res.TrackerResult)
}

func formatTrackerResult(res tests.TrackerResult) {
	fmt.Printf("Tracker supports IP spoofing: %t\n", res.SupportsIPSpoofing)
	fmt.Printf("Tracker supports optimized announce responses: %t\n", res.SupportsAnnouncingPeerNotInPeerList)
	fmt.Printf("Tracker supports optimized seeder announce responses: %t\n", res.SupportsOptimizedSeederResponse)

	fmt.Println()
	fmt.Println("Poke ran these tests:")
	for _, t := range res.Tests {
		if !t.Run {
			continue
		}
		fmt.Println()
		fmt.Printf("Test: %s\n", t.Name)
		if t.Result.Result != nil {
			fmt.Printf("Result: %v\n", t.Result.Result)
		}
		if t.Result.Err != nil {
			fmt.Printf("Error: %s\n", t.Result.Err)
		}
	}

	fmt.Println()
	fmt.Println("Poke did not run these tests:")
	for _, t := range res.Tests {
		if t.Run {
			continue
		}
		fmt.Printf("%s\t- %s\n", t.Name, t.NotRunReason)
	}
}

func runOtherTests(announceURI string) {
	err := tests.BasicHTTPSeederAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.CheckReturnedPeersHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidShortInfohashHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidLongInfohashHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidShortPeerIDHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidLongPeerIDHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidNegativeUploadedHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidNegativeDownloadedHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidNegativeLeftHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.InvalidEventHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}
}
