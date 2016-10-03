# poke
A tool to feature-test BitTorrent trackers via HTTP and UDP announces

# How to get it
A simple

    go get github.com/mrd0ll4r/poke/cmd/poke

should get and build it.

# How to use it
Run

    poke -a <announce URI> [-debug]

to test the tracker specified by `<announce URI>` via HTTP.
To use UDP, specify the UDP endpoing (e.g. `localhost:1234`) via the `-u` flag.
The `-u` flag has priority over the `-a` flag.

# License
MIT