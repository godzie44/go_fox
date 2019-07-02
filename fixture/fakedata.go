package fixture

import (
	"database/sql"
	"github.com/Pallinder/go-randomdata"
	"github.com/lib/pq"
	"log"
	"math/rand"
	"time"
)

func GenerateData(db *sql.DB)  {
	log.Println("Creating table fox_test")

	_, err := db.Query(`DROP TABLE IF EXISTS fox_test`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query(`CREATE TABLE fox_test (user_id BIGINT NOT NULL, ip_addr VARCHAR(15) NOT NULL, ts TIMESTAMP NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query(`CREATE INDEX user_id_idx ON fox_test USING hash(user_id)`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("fill by random data")

	var ipPool = make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		ipPool = append(ipPool, randomdata.IpV4Address())
	}

	var userIpList = make(map[int64]*[]string)

	for i := 0; i < 100; i++ {

		trans, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := trans.Prepare(pq.CopyIn("fox_test", "user_id", "ip_addr", "ts"))
		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()
		pqTime := pq.FormatTimestamp(t)

		for j := 0; j < 100000; j++ {
			userId := rand.Int63n(5000000)

			if _, ok := userIpList[userId]; !ok {
				newIpList := make([]string, 0, 3)
				userIpList[userId] = &newIpList
			}

			var userIpToInsert string
			ipList := userIpList[userId]

			if len(*ipList) > 2 {
				userIpToInsert = (*ipList)[rand.Intn(3)]
			} else {
				userIpToInsert = ipPool[rand.Intn(100)]
				*ipList = append(*ipList, userIpToInsert)
			}

			_, err = stmt.Exec(userId, userIpToInsert, pqTime)
			if err != nil {
				log.Fatal(err)
			}

		}

		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		err = stmt.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = trans.Commit()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Add 100 000 rows")

	}
}