package shorter

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

var validShortURL = []string{
	"0",
	"1",
	"a",
	"abcd",
	"edfzds",
	"012abcAC",
	"lYGhA16ahyf", //maxuint64
}
var invalidShortURL = []string{
	">",
	"<",
	"lYGhA16ahyf00", //overflow
	"abcd>",
}

var validOrigURL = []string{
	"http://google.com",
	"scheme://test.xx",
	"https://abc.xyz",
	"http://127.0.0.1",
}
var invalidOrigURL = []string{
	" /main.html",
	" http:www.example.com/main.html",
	" www.example.com/main.html",
	"hlo://\\",
	"google.com/xxx&111",
}

func getDefaultDB() *DefaultDB {
	db := DefaultDB{}
	db.Init()
	return &db
}
func TestGetID(t *testing.T) {

	for _, valid := range validShortURL {
		_, err := GetID(valid)
		if err != nil {
			t.Error("get id failed:", err)
			return
		}
	}
	for _, invalid := range invalidShortURL {
		id, err := GetID(invalid)
		if err == nil {
			t.Error("should fail of invalid short URL:", invalid, id)
			return
		}

	}

}
func BenchmarkShorterURLGene(b *testing.B) {
	for i := uint64(0); i < uint64(b.N); i++ {
		ShorterURLGene(i)
	}
}
func initDB() {
	db := getDefaultDB()
	InitWithDB(db)
}
func TestGetOrigURLByShortURL(t *testing.T) {
	_, err := GetOrigURLByShortURL("test")
	if err != nil {
		if !strings.Contains(err.Error(), "db not init") {
			t.Error("should return not init")
			return
		}
	} else {
		t.Error("should return err")
		return
	}
	initDB()
	for _, origURL := range validOrigURL {
		short, err := ShorterURL(origURL)
		if err != nil {
			t.Error("fail to get short:", err)
			return
		}
		tmp, err := GetOrigURLByShortURL(short)
		if err != nil {
			t.Error("fail to get orig by short:", err)
			return
		}
		if tmp != origURL {
			t.Error("get orig url not equal orig data")
			return
		}
	}
	for _, origURL := range invalidOrigURL {
		short, err := ShorterURL(origURL)
		if err == nil {

			t.Error("should return err:", origURL, short)
		}
	}
}

func TestShorterURLGene(t *testing.T) {
	for i := uint64(0); i < 100; i++ {
		short := ShorterURLGene(i)
		origID, err := GetID(short)
		if err != nil {
			t.Error("fail to getid from short:", err)
			return
		}
		if origID != i {
			t.Error("id not equal")
			return
		}
	}
	for i := uint64(math.MaxUint64 - 100); i != 0 && i <= math.MaxUint64; i++ {

		short := ShorterURLGene(i)
		origID, err := GetID(short)
		if err != nil {
			t.Error("fail to getid from short:", err)
			return
		}
		if origID != i {
			t.Error("id not equal")
			return
		}

	}
}

func TestMulShorterURLs(t *testing.T) {
	shortURLs, err := MulShorterURLs(validOrigURL)
	if err != nil {
		t.Error("fal to mul shorter urls:", err)
		return
	}
	for i, short := range shortURLs {
		origURL, err := GetOrigURLByShortURL(short)
		if err != nil {
			t.Error("fail to get orig by short:", err, short, origURL)
			return
		}
		if origURL != validOrigURL[i] {
			t.Error("get orig not equal:", err)
			return
		}
	}
}

func TestMulShorterURLsGene(t *testing.T) {
	var id []uint64
	for i := uint64(0); i < 100; i++ {
		id = append(id, i)
	}
	for i := uint64(math.MaxUint64 - 100); i != 0 && i <= math.MaxUint64; i++ {
		id = append(id, i)
	}
	fmt.Println("gene id finish：", len(id))
	shortURLs := MulShorterURLsGene(id)
	for i := uint64(0); i < 100; i++ {
		origID, err := GetID(shortURLs[i])
		if err != nil {
			t.Error("fail to get id:", shortURLs[i])
			return
		}

		if origID != i {
			t.Error("not equal:", origID, i)
			return
		}
	}
	for i := 100; i < len(id); i++ {
		origID, err := GetID(shortURLs[i])
		if err != nil {
			t.Error("fail get id:", shortURLs[i])
			return
		}

		if origID != id[i] {
			t.Error("not equal:", origID, id[i])
			return
		}
	}
}
