package torrent

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"os"
	"reflect"

	"github.com/marconi/byteslego/tracker"
	"github.com/zeebo/bencode"
)

type Torrent struct {
	Filename string
	Metainfo map[string]interface{}
	InfoHash string
	PeerId   string
}

/*
 * API
 */

func New(f string) *Torrent {
	metaInfo := readMetainfo(f)
	torrent := &Torrent{
		Filename: f,
		Metainfo: metaInfo,
		InfoHash: genInfoHash(metaInfo),
		PeerId:   genPeerId(),
	}
	return torrent
}

func (t *Torrent) GetTracker() *tracker.Tracker {
	announce := t.Metainfo["announce"].(string)
	tracker := tracker.New(announce)
	return tracker
}

func (t *Torrent) GetTrackers() []*tracker.Tracker {
	var trackers []*tracker.Tracker
	for _, a := range t.Metainfo["announce-list"].([]interface{}) {
		i := reflect.ValueOf(a).Index(0)
		url := i.Interface().(string)
		trackers = append(trackers, tracker.New(url))
	}
	return trackers
}

/*
 * Utilities
 */
func readMetainfo(f string) map[string]interface{} {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	dec := bencode.NewDecoder(r)

	var metaInfo interface{}
	if err := dec.Decode(&metaInfo); err != nil {
		panic(err)
	}
	return iterate(metaInfo).(map[string]interface{})
}

func iterate(ival interface{}) interface{} {
	var iterateOut interface{}
	tVal := reflect.ValueOf(ival)
	switch tVal.Kind() {
	case reflect.Map:
		mVal := tVal.Interface().(map[string]interface{})
		mapOut := make(map[string]interface{})
		for key, val := range mVal {
			switch reflect.ValueOf(val).Kind() {
			case reflect.Map, reflect.Slice:
				mapOut[key] = iterate(val)
			default:
				if key == "pieces" {
					mapOut[key] = "..."
				} else {
					mapOut[key] = val
				}
			}
		}
		iterateOut = mapOut
	case reflect.Slice:
		mVal := tVal.Interface().([]interface{})
		var sliceOut []interface{}
		for _, val := range mVal {
			switch reflect.ValueOf(val).Kind() {
			case reflect.Map, reflect.Slice:
				sliceOut = append(sliceOut, iterate(val))
			default:
				sliceOut = append(sliceOut, val)
			}
		}
		iterateOut = sliceOut
	}
	return iterateOut
}

func genInfoHash(metaInfo interface{}) string {
	buffer := new(bytes.Buffer)
	enc := bencode.NewEncoder(buffer)
	if err := enc.Encode(metaInfo); err != nil {
		panic(err)
	}
	hash := sha1.New()
	hash.Write(buffer.Bytes())
	return string(hash.Sum(nil))
}

func genPeerId() string {
	pid := new([20]byte)
	_, err := rand.Read(pid[:])
	if err != nil {
		panic(err)
	}
	return string(pid[:])
}
