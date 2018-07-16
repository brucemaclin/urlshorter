package shorter

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

//DB db func to handle shortURL info
type DB interface {

	//Add save shortURL and origURL with the origID
	Add(origURL string) (string, error)

	//GetOrigURLByShort get orig url will query sql by  shorturl
	GetOrigURLByShort(shortURL string) (origURL string, err error)

	//GetNextID will get a unused auto incr id.
	//used if shorter AutoGetNextID enabled
	GetNextID() (uint64, error)
	//Clone the storage if needed.For example,using mgo,you can clone the session with session.Clone
	//to avoid cocurrent access problems.
	//can return it self if not a problem
	Clone() DB

	//MulAdd add origURLs save in db and return shortURLs
	MulAdd(origURLs []string) ([]string, error)
}

//*DefaultDB will use memory to store infos of shortURL
type DefaultDB struct {
	sync.RWMutex
	nextID      uint64
	shortURLMap map[string]urlInfo
	origURLMap  map[string]urlInfo
}
type urlInfo struct {
	origURL    string
	id         uint64
	createTime time.Time
	shortURL   string
}

//Init inner map.should be call before all other func
func (db *DefaultDB) Init() {
	db.shortURLMap = make(map[string]urlInfo)
	db.origURLMap = make(map[string]urlInfo)
	db.nextID += 100 + uint64(rand.Int63n(1000))
}

//Clone return db self
func (db *DefaultDB) Clone() DB {
	return db
}

//Add save url info in map
func (db *DefaultDB) Add(origURL string) (string, error) {
	db.Lock()
	defer db.Unlock()
	origID := db.nextID
	db.nextID++
	shortURL := convert10To62(origID)
	var info = urlInfo{origURL: origURL, id: origID, createTime: time.Now(), shortURL: shortURL}
	db.shortURLMap[shortURL] = info
	db.origURLMap[origURL] = info
	return shortURL, nil
}

//GetNextID return nextID and increment nextID
func (db *DefaultDB) GetNextID() (uint64, error) {
	db.Lock()
	defer db.Unlock()
	res := db.nextID
	db.nextID++
	return res, nil
}

//GetOrigURLByShort return origURL by query map with short url
func (db *DefaultDB) GetOrigURLByShort(shortURL string) (string, error) {
	db.RLock()
	defer db.RUnlock()
	info, ok := db.shortURLMap[shortURL]
	if !ok {
		return "", errors.New("not exist this url")
	}
	return info.origURL, nil
}

func (db *DefaultDB) checkIfExistShortURL(origURL string) (string, bool) {
	db.RLock()
	defer db.RUnlock()
	info, ok := db.origURLMap[origURL]
	if !ok {
		return "", ok
	}
	return info.shortURL, ok
}

func (db *DefaultDB) MulAdd(origURLs []string) ([]string, error) {
	if len(origURLs) == 0 {
		return nil, errors.New("no orig url")
	}
	db.Lock()
	defer db.Unlock()
	var res []string
	start := db.nextID
	end := db.nextID + uint64(len(origURLs))
	if end < start {
		return nil, errors.New("id overflow")
	}
	for i := start; i < end; i++ {
		shortURL := convert10To62(i)
		res = append(res, shortURL)
		var info = urlInfo{origURL: origURLs[i-start], id: start, createTime: time.Now(), shortURL: shortURL}
		db.shortURLMap[shortURL] = info
		db.origURLMap[origURLs[i-start]] = info
	}
	return res, nil
}
