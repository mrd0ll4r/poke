# poke
A tool to feature-test BitTorrent trackers via HTTP

# How to get it
A simple

    go get github.com/mrd0ll4r/poke/cmd/poke

should get and build it.

# How to use it
Run

    poke -a <announce URI> [-debug]

to test the tracker specified by `<announce URI>` via HTTP.

# License
MIT