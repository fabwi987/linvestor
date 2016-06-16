package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbUser string
var dbPassword string
var dbServer string
var dbName string
var dbType string
var dbTable string
var connString string
var tableExist bool

//DbCreateConnectionString assembles the parameters from the config file to a connection string
func DbCreateConnectionString() {
	GetSQLConfigParameters()
	connString = dbUser + ":" + dbPassword + "@" + dbServer + "/" + dbName
	//log.Println("Connection string: " + connString)

}

//DbTableCreate creates the tabel the first time the application is used
func DbTableCreate() {
	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)
	//log.Println(dbTable)

	rows, err := db.Query("SHOW TABLES LIKE '" + dbTable + "'")
	if !rows.Next() {
		log.Println("Creating tabel")
		stmt, err := db.Prepare("CREATE TABLE " + dbTable + " (idstock1 int(11) NOT NULL AUTO_INCREMENT,symbol varchar(45) DEFAULT NULL,created varchar(45) DEFAULT NULL,buyprice float DEFAULT NULL,noshares float DEFAULT NULL,PRIMARY KEY (idstock1)) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8")
		res, err := stmt.Exec()
		Perror(err)
		log.Println(res.RowsAffected())
	} else {
		log.Println("Will not create Tabel")
	}

	tableExist = true

}

//DBTestConnection is a simple method to test connection to DB
func DBTestConnection() {
	DbCreateConnectionString()

	//connString = "b06fa04f33f026:a480ce77@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_4499eae5e6a7bd4" //FUNKAR!!

	db, err := sql.Open(dbType, connString)
	Perror(err)

	var noRows int
	rows, err := db.Query("SELECT COUNT(*) FROM" + dbTable)
	log.Println(err)
	rows.Next()
	rows.Scan(&noRows)
	log.Println(noRows)
	db.Close()
}

//DbInsertSQL inserts a new stock into an sql server database
func DbInsertSQL(_stock StockDataSaveFormat) (string, error) {
	//if !tableExist {
	//DbTableCreate()
	//}

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	stmt, err := db.Prepare("INSERT " + dbTable + " SET symbol=?,created=?,buyprice=?,noshares=?")
	Perror(err)

	res, err := stmt.Exec(_stock.Symbol, _stock.Created, _stock.BuyPrice, _stock.NumberOfShares)
	Perror(err)
	nid, err := res.LastInsertId()
	log.Println(nid)

	db.Close()

	log.Println("Insert i DB")

	return "Insert successfull", nil
}

//DBQuerySQL selects all stocks from the database
func DBQuerySQL() ([]StockDataSaveFormat, error) {
	if !tableExist {
		DbTableCreate()
	}

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	var noRows int
	rows, err := db.Query("SELECT COUNT(*) FROM " + dbTable)
	Perror(err)
	for rows.Next() {
		err = rows.Scan(&noRows)
	}
	Perror(err)

	var newStocks = make([]StockDataSaveFormat, noRows)

	rows, err = db.Query("SELECT * FROM " + dbTable)
	Perror(err)

	var i int
	for rows.Next() {

		var tempStock StockDataSaveFormat
		var id int
		err = rows.Scan(&id, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares)
		Perror(err)
		newStocks[i] = tempStock
		//log.Println(newStocks[i].Symbol)
		i++
	}

	log.Println("Hämtar data från DB")
	return newStocks, nil
}

//GetSQLConfigParameters sets the parameters to connect to database
func GetSQLConfigParameters() {

	/**
	viper.SetConfigType("yaml")
	viper.SetConfigName("db_config")
	viper.AddConfigPath("/configs/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println(err)
		log.Println("No configuration file loaded - using defaults")
	}


	dbUser = viper.GetString("db_user")
	dbPassword = viper.GetString("db_password")
	dbServer = viper.GetString("db_server")
	dbName = viper.GetString("db_name")
	dbType = viper.GetString("db_type")
	dbTable = viper.GetString("db_table")**/

	dbUser = "b0e21cd256f9bf"
	dbPassword = "48122e8a"
	dbServer = "tcp(us-cdbr-iron-east-04.cleardb.net:3306)"
	dbName = "heroku_c52163db78254a1"
	dbType = "mysql"
	dbTable = "stock123"

}
