// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package tmpbbs

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en_US": &dictionary{index: en_USIndex, data: en_USData},
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"%d replies":             0,
	"%dd ago":                3,
	"%dh ago":                4,
	"%dm ago":                2,
	"Author#tripcode-secret": 6,
	"Close":                  11,
	"Insert emoji using shortcode between colons (:mushroom: becomes 🍄).": 8,
	"Markdown is supported.": 7,
	"Reply":                  9,
	"Title":                  5,
	"URL QR Code":            10,
	"page %d":                1,
}

var en_USIndex = []uint32{ // 13 elements
	0x00000000, 0x00000024, 0x0000002f, 0x0000003a,
	0x00000045, 0x00000050, 0x00000056, 0x0000006d,
	0x00000084, 0x000000cb, 0x000000d1, 0x000000dd,
	0x000000e3,
} // Size: 76 bytes

const en_USData string = "" + // Size: 227 bytes
	"\x14\x01\x81\x01\x00=\x01\x0c\x02%[1]d reply\x00\x0e\x02%[1]d replies" +
	"\x02page %[1]d\x02%[1]dm ago\x02%[1]dd ago\x02%[1]dh ago\x02Title\x02Aut" +
	"hor#tripcode-secret\x02Markdown is supported.\x02Insert emoji using shor" +
	"tcode between colons (:mushroom: becomes 🍄).\x02Reply\x02URL QR Code\x02" +
	"Close"

	// Total table size 303 bytes (0KiB); checksum: FF5AB438
