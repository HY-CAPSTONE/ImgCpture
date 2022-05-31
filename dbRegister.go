package main

import (
	"database/sql"
	"fmt"
	"log"
)

/*
- 특정이미지 특정이름으로 바꿔서 디렉토리에 저장 하는 함수
- 해당 디렉토리 이미지를 DB에 등록하는 함수
- URL을 받아서 이미지를 노출하는 함수
*/

func insert_to_db(dbInfo dbConn, imgDir string) {
	qurey := "INSERT INTO Gallery (PLANT_ID, SAVE_PATH) VALUES ( 99, " + imgDir + " );"

	driverName, sataSourceName := dbInfo.dbAddress()
	// Create the database handle, confirm driver is present
	db, err := sql.Open(driverName, sataSourceName)
	defer db.Close()
	if err != nil || db.Ping() != nil {
		panic(err.Error())
	}

	// Connect and check the server version
	var PLANT_ID, SAVE_DATE, SAVE_PATH string
	rows, err := db.Query(qurey)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&PLANT_ID, &SAVE_DATE, &SAVE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(PLANT_ID, SAVE_DATE, SAVE_PATH)
	}
}
