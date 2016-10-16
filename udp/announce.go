package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"

	"github.com/mrd0ll4r/poke"
)

func init() {
	var transid = rand.Uint32()
	tid = &transid
}

// ErrInvalidAddress indicates an invalid address was used to create a client.
var ErrInvalidAddress = errors.New("invalid address")

var tid *uint32

// Client is a UDP client.
type Client struct {
	addr         string
	conn         net.Conn
	connectionID uint64
	autoConnect  bool
}

var _ poke.Announcer = &Client{}

// SetAutoConnect enables or disables the automatic creation of connection IDs.
func (c *Client) SetAutoConnect(to bool) {
	c.autoConnect = to
}

// SetConnectionID sets the connection ID to use for future requests.
//
// This is only used if AutoConnect is set to false.
func (c *Client) SetConnectionID(to uint64) {
	c.connectionID = to
}

// NewClient creates a new client for the given tracker address.
//
// The Client will automatically make connect requests for every announce.
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		addr:        addr,
		conn:        conn,
		autoConnect: true,
	}, nil
}

// ManualConnect performs a connect request and returns the connection ID.
func (c *Client) ManualConnect() (uint64, error) {
	return c.manualConnect()
}

func (c *Client) manualConnect() (uint64, error) {
	transactionID := atomic.AddUint32(tid, 1)

	buf := make([]byte, 16)

	buf[2] = 0x04
	buf[3] = 0x17
	buf[4] = 0x27
	buf[5] = 0x10
	buf[6] = 0x19
	buf[7] = 0x80

	binary.BigEndian.PutUint32(buf[12:16], transactionID)

	n, err := c.conn.Write(buf)
	if err != nil {
		return 0, err
	}

	if n != 16 {
		return 0, errors.New("connect: Did not send 16 bytes")
	}

	buf = make([]byte, 64)
	n, err = c.conn.Read(buf)
	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") {
			return 0, errors.New("connect: I/O timeout on receive")
		}
		return 0, fmt.Errorf("connect: %s", err)
	}

	if n != 16 {
		return 0, errors.New("connect: Did not receive 16 bytes")
	}

	b := buf[:n]

	action := binary.BigEndian.Uint32(b[:4])
	if action != 0 {
		return 0, errors.New("connect: action != 0")
	}

	transID := binary.BigEndian.Uint32(b[4:8])
	if transID != transactionID {
		return 0, errors.New("connect: transaction IDs do not match")
	}

	connID := binary.BigEndian.Uint64(b[8:16])
	return connID, nil
}

func prepareAnnounce(req poke.AnnounceRequest, connID uint64, transactionID uint32) ([]byte, error) {
	bbuf := bytes.NewBuffer(nil)

	// Connection ID
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, connID)
	_, err := bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Action
	_, err = bbuf.Write([]byte{0, 0, 0, 0x01})
	if err != nil {
		return nil, err
	}

	// Transaction ID
	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, transactionID)
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Infohash
	_, err = bbuf.Write(req.InfoHash)
	if err != nil {
		return nil, err
	}

	// Peer ID
	_, err = bbuf.WriteString(req.Peer.ID)
	if err != nil {
		return nil, err
	}

	// Downloaded
	buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(req.Downloaded))
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Left
	buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(req.Left))
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Uploaded
	buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(req.Uploaded))
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Event
	_, err = bbuf.Write([]byte{0, 0, 0, 0x02})
	if err != nil {
		return nil, err
	}

	// IP Address
	_, err = bbuf.Write(req.Peer.IP)
	if err != nil {
		return nil, err
	}

	// "Key"
	_, err = bbuf.Write([]byte{0, 0, 0, 0})
	if err != nil {
		return nil, err
	}

	// Numwant
	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(req.Numwant))
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	// Port
	buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, req.Port)
	_, err = bbuf.Write(buf)
	if err != nil {
		return nil, err
	}

	return bbuf.Bytes(), nil
}

// Announce performs an announce for this client.
// If autoConnect is enabled, the client will first perform a connect request to
// obtain a connection ID.
//
// This implements poke.Announcer for UDP clients.
func (c *Client) Announce(req poke.AnnounceRequest) (poke.OptionalAnnounceResponse, error) {
	if ip := req.IP.To4(); ip != nil {
		req.IP = ip
	}

	if poke.Debug {
		log.Printf("Announcing: %+v", req)
	}

	if c.autoConnect {
		connID, err := c.manualConnect()
		if err != nil {
			return nil, err
		}
		c.connectionID = connID
	}

	transactionID := atomic.AddUint32(tid, 1)
	toReturn := poke.AnnounceResponse{
		Peers: make([]poke.Peer, 0),
	}
	packet, err := prepareAnnounce(req, c.connectionID, transactionID)
	if err != nil {
		return nil, err
	}

	// Prepare a receive buffer.
	buf := make([]byte, 1024)

	// Send announce.
	_, err = c.conn.Write(packet)
	if err != nil {
		return nil, err
	}

	// Receive response.
	n, err := c.conn.Read(buf)
	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") {
			return nil, errors.New("announce: I/O timeout on receive")
		}
		return nil, fmt.Errorf("announce: %s", err)
	}
	if n < 20 {
		return nil, errors.New("announce: Did not receive at least 20 bytes")
	}

	// Check transaction ID.
	transID := binary.BigEndian.Uint32(buf[4:8])
	if transID != transactionID {
		return nil, errors.New("announce: transaction IDs do not match")
	}

	// Parse action.
	action := binary.BigEndian.Uint32(buf[:4])
	if action != 1 {
		if action == 3 {
			errVal := string(buf[8:n])
			return poke.ErrorResponse(errVal), nil
			//return nil, fmt.Errorf("announce: tracker responded with error: %s", errVal)
		}
		return nil, errors.New("announce: tracker responded with action != 1")
	}

	toReturn.Interval = int(binary.BigEndian.Uint32(buf[8:12]))
	toReturn.Incomplete = int(binary.BigEndian.Uint32(buf[12:16]))
	toReturn.Complete = int(binary.BigEndian.Uint32(buf[16:20]))

	if (n-20)%6 != 0 {
		return nil, fmt.Errorf("announce: unexpected announce response length: %d", n)
	}

	numPeers := (n - 20) / 6
	for i := 0; i < numPeers; i++ {
		toReturn.Peers = append(toReturn.Peers,
			poke.Peer{
				IP:   net.IP(buf[20+6*i : 24+6*i]),
				Port: binary.BigEndian.Uint16(buf[24+6*i : 26+6*i]),
			})
	}

	if poke.Debug {
		log.Printf("Got announce response: %+v", toReturn)
	}

	return toReturn, nil
}
