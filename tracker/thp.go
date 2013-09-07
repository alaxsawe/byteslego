/**
 * Tracker HTTP Protocol
 */

package tracker

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zeebo/bencode"
)

const (
	STARTED   = "started"
	STOPPED   = "stopped"
	COMPLETED = "completed"
)

type Tracker struct {
	Announce string
}

type Request struct {
	InfoHash   string
	PeerId     string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Ip         string
	NumWant    int
	Event      string // started, stopped, completed
	Announce   string
	Url        string
}

func NewRequest(announce, ih, pid string, port int, event string) (*Request, error) {
	r := &Request{
		InfoHash: ih,
		PeerId:   pid,
		Port:     port,
		Ip:       "112.210.47.38",
		NumWant:  5,
		Event:    event,
		Announce: announce,
	}

	u, err := url.Parse(r.Announce)
	if err != nil {
		return nil, err
	}

	uq := u.Query()
	uq.Add("info_hash", r.InfoHash)
	uq.Add("peer_id", r.PeerId)
	uq.Add("port", strconv.Itoa(r.Port))
	uq.Add("ip", r.Ip)
	uq.Add("numwant", strconv.Itoa(r.NumWant))
	uq.Add("uploaded", strconv.Itoa(r.Uploaded))
	uq.Add("downloaded", strconv.Itoa(r.Downloaded))
	uq.Add("left", strconv.Itoa(r.Left))
	uq.Add("event", r.Event)

	// attach query string
	u.RawQuery = uq.Encode()

	r.Url = u.String()
	return r, nil
}

func New(announce string) *Tracker {
	return &Tracker{Announce: announce}
}

func (t *Tracker) GetPeers(r *Request) ([]map[string]string, error) {
	res, err := http.Get(r.Url)
	if err != nil {
		return nil, err
	}

	dec := bencode.NewDecoder(res.Body)
	defer res.Body.Close()
	var resp interface{}
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}

	fmt.Println(resp)
	return nil, nil
}
