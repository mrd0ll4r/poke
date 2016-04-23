package main

import (
	"log"

	"github.com/mrd0ll4r/poke"
	"github.com/mrd0ll4r/poke/tests"
	"flag"
)

func init() {
	flag.StringVar(&announceURI,"a","http://localhost:6882/announce","the announce URI")
	flag.BoolVar(&debug,"debug",false,"debug mode")
}

var (
	announceURI string
	debug bool
)

func main() {
	flag.Parse()

	poke.Debug = debug

	err := tests.BasicHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.BasicHTTPNonCompactAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	err = tests.BasicHTTPSeederAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	}

	pruning,err := tests.TrackerSupportsAnnouncingPeerNotInPeerListHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	} else {
		if pruning {
			log.Println("tracker supports leaving the announcing peer out of the peer list")
		} else {
			log.Println("tracker does not support leaving the announcing peer out of the peer list")
		}
	}

	compact, err := tests.TrackerSupportsCompactHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	} else {
		if compact {
			log.Println("tracker supports compact HTTP announces")
		} else {
			log.Println("tracker does not support compact HTTP announces")
		}
	}

	nonCompact, err := tests.TrackerSupportsNonCompactHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	} else {
		if nonCompact {
			log.Println("tracker supports non-compact HTTP announces")
		} else {
			log.Println("tracker does not support non-compact HTTP announces")
		}
	}

	optimized, err := tests.TrackerSupportsOptimizedSeederHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	} else {
		if optimized {
			log.Println("tracker supports optimized responses for seeder announces (does not return other seeders)")
		} else {
			log.Println("tracker does not support optimized responses for seeder announces (returns other seeders)")
		}
	}

	spoofing, err := tests.TrackerSupportsIPSpoofingHTTPAnnounce(announceURI)
	if err != nil {
		log.Println(err)
	} else {
		if spoofing {
			log.Println("tracker supports IP spoofing")
		} else {
			log.Println("tracker does not support IP spoofing")
		}
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
