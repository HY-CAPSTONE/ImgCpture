package module

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

/*
- 특정이미지 특정이름으로 바꿔서 디렉토리에 저장 하는 함수
- 해당 디렉토리 이미지를 DB에 등록하는 함수
- URL을 받아서 이미지를 노출하는 함수
*/

type DBConn struct {
	DriverName  string
	UserName    string
	Password    string
	Ip          string
	Port        string
	TargetTable string
}

func (db DBConn) DBAddress() (string, string) {
	dataSourceName := db.UserName + ":" + db.Password + "@tcp(" + db.Ip + ":" + db.Port + ")/" + db.TargetTable
	return db.DriverName, dataSourceName
}

func Insert_to_db(sendImg *chan bool, dbInfo DBConn, imgDir string) {
	qurey := "INSERT INTO Gallery (PLANT_ID, SAVE_PATH) VALUES ( 99, \"" + imgDir + "\" );"

	driverName, dataSourceName := dbInfo.DBAddress()
	// Create the database handle, confirm driver is present
	db, err := sql.Open(driverName, dataSourceName)
	defer db.Close()
	if err != nil || db.Ping() != nil {
		panic(err.Error())
	}
	println("SQL is Connected")

	// Connect and check the server version
	//var PLANT_ID, SAVE_DATE, SAVE_PATH string
	println(qurey)
	db.Query(qurey)
	if err != nil {
		log.Fatal(err)
	}
	*sendImg <- true
}

func FindRecentImg(ownerName string, dbInfo DBConn) string {
	qurey := "SELECT SAVE_PATH FROM Gallery WHERE PLANT_ID=" + ownerName + " ORDER BY SAVE_DATE DESC LIMIT 1;"
	driverName, dataSourceName := dbInfo.DBAddress()
	// Create the database handle, confirm driver is present
	db, err := sql.Open(driverName, dataSourceName)
	defer db.Close()
	if err != nil || db.Ping() != nil {
		panic(err.Error())
	}
	println("SQL is Connected")

	// Connect and check the server version
	//var PLANT_ID, SAVE_DATE, SAVE_PATH string
	var recentDir string
	err = db.QueryRow(qurey).Scan(&recentDir)
	if err != nil {
		log.Fatal(err)
	}
	return recentDir
}
