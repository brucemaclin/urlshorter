package shorter

import (
	"sync/atomic"
	"errors"
	"net/url"
)

type shortInfo struct {
	db DB
	dbFlag int32
}
var inner shortInfo
func InitWithDB(db DB) {
	inner.db = db
	atomic.AddInt32(&inner.dbFlag,1)
	return
}

//ShorterURL can use after InitWithDB
//it will GetNextID and then use db to save infos
//of url.
//ret: shortURL and error if db not init or error occurs when use db
func ShorterURL(origURL string) (string,error) {
	if origURL == ""  {
		return "",errors.New("origURL should not be blank")
	}

	uri,err := url.Parse(origURL)
	if err != nil {
		return "",err
	}

	if uri.Host == ""  {
		return "",errors.New("uri has no host")
	}
	val := atomic.LoadInt32(&inner.dbFlag)
	if val == 0 {
		return "",errors.New("db not init")
	}
	id ,err := inner.db.GetNextID()
	if err != nil {
		return "",err
	}
	short := convert10To62(id)
	err = inner.db.Add(id,short,origURL)
	if err != nil {
		return "",err
	}
	return short,err
}

//ShorterURLGene only gene a shorter url with id
//use 62 decimal
func ShorterURLGene(id uint64) string {
	return convert10To62(id)
}

//GetID return 10 decimal id if shortURL valid
//you can use id to query sql or use this func in db.GetOrigURLByShort
func GetID(shortURL string) (uint64,error){
	return convert62To10(shortURL)
}
//GetOrigURLByShortURL get origURL by short
//db should init before.
//will call db.GetOrigURLByShortURL to query db to get origURL
//ret: origURL and error
func GetOrigURLByShortURL(short string) (string,error) {
	val := atomic.LoadInt32(&inner.dbFlag)
	if val == 0 {
		return "",errors.New("db not init")
	}
	return inner.db.GetOrigURLByShort(short)
}