package search

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type SearchEngine struct {
	Db                *sql.DB
	//save in cache only duplicated users, cause for not duplicated its can change in future
	cacheOfDuplicates map[string]interface{}
	sync.Mutex
}

func NewSearchEngine(db *sql.DB) *SearchEngine {
	return &SearchEngine{
		Db:db,
		cacheOfDuplicates: make(map[string]interface{}),
	}
}

type SearchResponse struct {
	Dupes bool
}

func (se *SearchEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	firstUserId, err := strconv.Atoi(vars["id1"])
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	secondUserId, err := strconv.Atoi(vars["id2"])
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	isDuplicated, err := se.usersIsDuplicated(firstUserId, secondUserId)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	s := SearchResponse{isDuplicated}
	res, err := json.Marshal(s)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (se *SearchEngine) usersIsDuplicated(userId1 int, userId2 int) (bool, error) {
	cacheKey := generateCacheKey(userId1, userId2)

	se.Lock()
	_, existsInCache := se.cacheOfDuplicates[cacheKey]
	se.Unlock()

	if existsInCache {
		log.Println("Hit in cache for ids ", userId1, userId2)
		return true, nil
	}

	isDuplicated, err := se.isDuplicatedWithoutCache(userId1, userId2)
	if err != nil {
		return false, err
	}

	if isDuplicated {
		se.Lock()
		se.cacheOfDuplicates[cacheKey] = struct{}{}
		se.Unlock()
	}

	return isDuplicated, nil
}

func generateCacheKey(userId1 int, userId2 int) string {
	return fmt.Sprintf("%d:%d", userId1, userId2)
}

func (se *SearchEngine) isDuplicatedWithoutCache(userId1 int, userId2 int) (bool, error) {
	user1Info := &UserInfo{}
	err := se.Db.QueryRow("select user_id, array_agg(ip_addr) from fox_test where user_id = $1 GROUP BY user_id", userId1).Scan(&user1Info.userId, pq.Array(&user1Info.ipAddressList))
	if err != nil {
		return false, err
	}

	user2Info := &UserInfo{}
	err = se.Db.QueryRow("select user_id, array_agg(ip_addr) from fox_test where user_id = $1 GROUP BY user_id", userId2).Scan(&user2Info.userId, pq.Array(&user2Info.ipAddressList))
	if err != nil {
		return false, err
	}

	log.Println(user1Info)
	log.Println(user2Info)

	return user1Info.isDuplicatedWith(user2Info), nil
}
