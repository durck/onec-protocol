package specification_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type pair struct {
	PayloadHex string `json:"payload_hex"`
	WireHex    string `json:"wire_hex"`
}
type message struct {
	HTTPText   string `json:"http_text"`
	HTTPLength int    `json:"http_length"`
	Target     string `json:"target"`
	Counter    uint32 `json:"counter"`
	RPCHex     string `json:"rpc_hex"`
	RPCLength  int    `json:"rpc_length"`
	WireHex    string `json:"wire_hex"`
}
type corpus struct {
	SchemaVersion int    `json:"schema_version"`
	Profile       string `json:"profile"`
	Origin        string `json:"origin"`
	NetworkUsed   bool   `json:"network_used"`
	Frames        []pair `json:"frames"`
	GUIDs         []struct {
		Value   string `json:"value"`
		WireHex string `json:"wire_hex"`
	} `json:"guids"`
	Integers []struct {
		Value   uint32 `json:"value"`
		WireHex string `json:"wire_hex"`
	} `json:"integers"`
	ByteAtoms []struct {
		DataHex string `json:"data_hex"`
		WireHex string `json:"wire_hex"`
	} `json:"byte_atoms"`
	OrdinaryList struct {
		Values  []uint32 `json:"values"`
		WireHex string   `json:"wire_hex"`
	} `json:"ordinary_list"`
	EmptyBlob   string `json:"empty_security_blob_hex"`
	TextStrings []struct {
		Value   string `json:"value"`
		WireHex string `json:"wire_hex"`
	} `json:"text_strings"`
	Introduction struct {
		Name    string `json:"name"`
		WireHex string `json:"wire_hex"`
	} `json:"introduction"`
	Files []struct {
		Length        int    `json:"length"`
		FillByte      byte   `json:"fill_byte"`
		ChunkLengths  []int  `json:"chunk_lengths"`
		EncodedLength int    `json:"encoded_length"`
		Hash          string `json:"encoded_sha256"`
	} `json:"http_files"`
	Request      message `json:"http_request"`
	Continuation message `json:"http_continuation"`
	Reply        message `json:"http_reply"`
}

func load(t *testing.T) corpus {
	t.Helper()
	b, err := os.ReadFile("vectors/vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var c corpus
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	if c.SchemaVersion != 1 || c.Profile != "native-18" || c.Origin != "synthetic" || c.NetworkUsed {
		t.Fatal("unexpected corpus scope")
	}
	if len(c.Frames) != 7 || len(c.GUIDs) != 4 || len(c.Integers) != 9 || len(c.ByteAtoms) != 5 || len(c.Files) != 5 || len(c.TextStrings) != 3 {
		t.Fatal("incomplete corpus")
	}
	return c
}

func unhex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func equalHex(t *testing.T, b []byte, expected string) {
	t.Helper()
	if !bytes.Equal(b, unhex(t, expected)) {
		t.Fatalf("encoding mismatch: got %x, want %s", b, expected)
	}
}

// These helpers implement only the encodings used by the published examples.
func integer(n uint32) []byte {
	if n < 10 {
		return []byte{byte(n + 1)}
	}
	if n < 256 {
		return []byte{0x0b, byte(n)}
	}
	if n < 65536 {
		return []byte{0x0d, byte(n), byte(n >> 8)}
	}
	return []byte{0x0f, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)}
}
func raw(b []byte) []byte {
	if len(b) == 0 {
		return []byte{1}
	}
	if len(b) == 1 {
		return []byte{0x0b, b[0]}
	}
	head := []byte{0x1a, byte(len(b))}
	if len(b) > 255 {
		head = []byte{0x1b, byte(len(b)), byte(len(b) >> 8)}
	}
	return append(head, b...)
}
func list(items ...[]byte) []byte {
	out := []byte{0x40}
	for i, b := range items {
		item := append([]byte(nil), b...)
		if i > 0 {
			item[0] |= 0x80
		}
		out = append(out, item...)
	}
	return append(out, 0x20)
}
func guid(t *testing.T, s string) []byte {
	t.Helper()
	b := unhex(t, strings.ReplaceAll(s, "-", ""))
	if len(b) != 16 {
		t.Fatal("bad GUID length")
	}
	return []byte{b[3], b[2], b[1], b[0], b[5], b[4], b[7], b[6], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}
}
func guidAtom(t *testing.T, s string) []byte { return append([]byte{0x15}, guid(t, s)...) }
func frame(b []byte) []byte {
	tokens := [][]byte{{0x65, 0x52, 0xb1, 0xa5}, {0x66, 0x53, 0xb2, 0xa6}, {0x21, 0x47, 0x54, 0xb3}, {0x22, 0x48, 0x55, 0xb5}}
	var out []byte
	for len(b) > 0 {
		matched := false
		for _, token := range tokens {
			if bytes.HasPrefix(b, token) {
				out = append(out, tokens[0]...)
				out = append(out, token...)
				b = b[4:]
				matched = true
				break
			}
		}
		if !matched {
			out = append(out, b[0])
			b = b[1:]
		}
	}
	return append(out, tokens[1]...)
}
func file(t *testing.T, b []byte) ([]byte, []int) {
	chunks := []byte{0x40}
	sizes := []int{}
	for len(b) >= 16384 {
		chunks = append(chunks, raw(b[:16384])...)
		sizes = append(sizes, 16384)
		b = b[16384:]
	}
	chunks = append(chunks, raw(b)...)
	sizes = append(sizes, len(b))
	chunks = append(chunks, 0x20)
	return list(raw([]byte("#")), guidAtom(t, "bea56fd3-3695-49fc-bbc5-c8bc15f33d41"), integer(1), integer(0), integer(0), integer(0), chunks), sizes
}
func emptyBlob(t *testing.T) []byte {
	// The inner length/data pair has no ordinary-list separator.
	return list(raw([]byte("#")), guidAtom(t, "671507fd-50a9-4b63-b70e-58b3d364f48f"), []byte{0x40, 1, 1, 0x20})
}

func TestScalarAndFramingVectors(t *testing.T) {
	c := load(t)
	for i, v := range c.Frames {
		t.Run(fmt.Sprintf("frame-%d", i), func(t *testing.T) { equalHex(t, frame(unhex(t, v.PayloadHex)), v.WireHex) })
	}
	for _, v := range c.GUIDs {
		equalHex(t, guid(t, v.Value), v.WireHex)
	}
	for _, v := range c.Integers {
		equalHex(t, integer(v.Value), v.WireHex)
	}
	for _, v := range c.ByteAtoms {
		equalHex(t, raw(unhex(t, v.DataHex)), v.WireHex)
	}
	var items [][]byte
	for _, n := range c.OrdinaryList.Values {
		items = append(items, integer(n))
	}
	equalHex(t, list(items...), c.OrdinaryList.WireHex)
	equalHex(t, emptyBlob(t), c.EmptyBlob)
	for _, v := range c.TextStrings {
		equalHex(t, []byte("\xef\xbb\xbf{\"S\",\""+strings.ReplaceAll(v.Value, "\"", "\"\"")+"\"}"), v.WireHex)
	}
	intro := []byte("V8IntroTicketReq")
	intro = binary.LittleEndian.AppendUint16(intro, uint16(len(c.Introduction.Name)))
	intro = append(intro, []byte(c.Introduction.Name)...)
	intro = append(intro, 0)
	equalHex(t, intro, c.Introduction.WireHex)
}

func TestHTTPFileBoundaries(t *testing.T) {
	for _, v := range load(t).Files {
		t.Run(strconv.Itoa(v.Length), func(t *testing.T) {
			b, sizes := file(t, bytes.Repeat([]byte{v.FillByte}, v.Length))
			if !reflect.DeepEqual(sizes, v.ChunkLengths) || len(b) != v.EncodedLength {
				t.Fatal("file length/chunk boundary mismatch")
			}
			sum := sha256.Sum256(b)
			equalHex(t, sum[:], v.Hash)
		})
	}
}

func TestHTTPEnvelopes(t *testing.T) {
	c := load(t)
	check := func(t *testing.T, b []byte, m message) {
		t.Helper()
		equalHex(t, b, m.RPCHex)
		equalHex(t, frame(b), m.WireHex)
		if len(b) != m.RPCLength {
			t.Fatal("RPC length mismatch")
		}
	}
	for _, m := range []message{c.Request, c.Reply} {
		parts := strings.SplitN(m.HTTPText, "\r\n\r\n", 2)
		if len(parts) != 2 || len(m.HTTPText) != m.HTTPLength {
			t.Fatal("HTTP message boundary mismatch")
		}
		want := "Content-Length: " + strconv.Itoa(len(parts[1]))
		if !strings.HasSuffix(parts[0], "\r\n"+want) {
			t.Fatal("Content-Length mismatch")
		}
	}
	for _, continued := range []bool{false, true} {
		m := c.Request
		if continued {
			m = c.Continuation
		}
		f, _ := file(t, []byte(c.Request.HTTPText))
		args := [][]byte{integer(0), guidAtom(t, c.Request.Target), integer(m.Counter)}
		if continued {
			args = append(args, integer(0))
		}
		args = append(args, integer(4), guidAtom(t, "bc15bd01-10bf-413c-a856-ddc907fcd123"), integer(0), integer(1), integer(3), integer(0), f, integer(2), emptyBlob(t), integer(2), list(raw([]byte("B")), integer(0)))
		check(t, list(args...), m)
	}
	f, _ := file(t, []byte(c.Reply.HTTPText))
	reply := list(integer(1), integer(c.Reply.Counter), integer(1), f, integer(3), integer(0), integer(2), emptyBlob(t), integer(2), list(raw([]byte("B")), integer(1)), integer(0))
	check(t, append(reply, 0x83), c.Reply)
}

func TestLocalDocumentLinks(t *testing.T) {
	link := regexp.MustCompile(`\]\(([^)]+)\)`)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range link.FindAllSubmatch(b, -1) {
			target := strings.SplitN(string(match[1]), "#", 2)[0]
			if target == "" || strings.Contains(target, "://") {
				continue
			}
			if filepath.IsAbs(target) {
				t.Errorf("%s: absolute link %s", path, target)
				continue
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(path), filepath.FromSlash(target))); err != nil {
				t.Errorf("%s: %s: %v", path, target, err)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
