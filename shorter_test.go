package shorter

import (
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
