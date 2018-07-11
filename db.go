package  shorter

import (
	"time"
	"sync"
	"errors"
)

//DB db func to handle shorturl info
type DB interface {

	//Add save shorturl and origurl with the origID
	Add(origID uint64,shortURL string, origURL string) error

	//GetOrigURLByShort get orig url will query sql by  shorturl
	GetOrigURLByShort(shortURL string) (origURL string,err error)

	//GetNextID will get a unused auto incr id.
	//used if shorter AutoGetNextID enabled
	GetNextID() (uint64,error)
	//Clone the storage if needed.For example,using mgo,you can clone the session with session.Clone
	//to avoid cocurrent access problems.
	//can return iteself if not a problem
	Clone() DB
}
//DefaultDB will use memory to store infos of shortURL
type DefaultDB struct {
	sync.RWMutex
	nextID uint64
	shortURLMap map[string]urlInfo
	origURLMap map[string]urlInfo
}
type urlInfo struct {
	origURL string
	id uint64
	createTime time.Time
	shortURL string

}
//Init inner map.should be call before all other func
func (db *DefaultDB) Init() {
	db.shortURLMap = make(map[string]urlInfo)
	db.origURLMap = make(map[string]urlInfo)
}
func (db *DefaultDB) Clone() DB {
	return db
}
func (db *DefaultDB) Add(origID uint64,shortURL string, origURL string) error {
	db.Lock()
	defer db.Unlock()
	var info = urlInfo{origURL:origURL,id:origID,createTime:time.Now(),shortURL:shortURL,}
	db.shortURLMap[shortURL] = info
	db.origURLMap[origURL] = info
	return nil
}
func (db *DefaultDB) GetNextID() (uint64,error) {
	db.Lock()
	defer db.Unlock()
	res := db.nextID
	db.nextID++
	return res,nil
}

func (db *DefaultDB) GetOrigURLByShort(shortURL string) (string,error ) {
	db.RLock()
	defer db.RUnlock()
	info,ok := db.shortURLMap[shortURL]
	if  !ok {
		return "",errors.New("not exist this url")
	}
	return info.origURL,nil
}

func (db *DefaultDB) checkIfExistShortURL(origURL string) (string,bool) {
	db.RLock()
	defer db.RUnlock()
	info,ok := db.origURLMap[origURL]
	if !ok {
		return "",ok
	}
	return info.shortURL,ok
}