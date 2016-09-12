package controllers

import (
	"database/sql"
	"log"

	"github.com/fabwi987/linvestor/models"
	_ "github.com/go-sql-driver/mysql"
)

var dbUser string
var dbPassword string
var dbServer string
var dbName string
var dbType string
var connString string
var tableExist bool

//DbCreateConnectionString assembles the parameters from the config file to a connection string
func DbCreateConnectionString() {
	GetSQLConfigParameters()
	connString = dbUser + ":" + dbPassword + "@" + dbServer + "/" + dbName
	//log.Println("Connection string: " + connString)

}

//DbTableCreate creates the tabel the first time the application is used
func DbTableCreate(dbtable string) {
	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)
	//log.Println(dbTable)

	rows, err := db.Query("SHOW TABLES LIKE '" + dbtable + "'")
	if !rows.Next() {
		log.Println("Creating tabel")
		stmt, err := db.Prepare("CREATE TABLE " + dbtable + " (idstock1 int(11) NOT NULL AUTO_INCREMENT,symbol varchar(45) DEFAULT NULL,created varchar(45) DEFAULT NULL,buyprice float DEFAULT NULL,noshares float DEFAULT NULL,PRIMARY KEY (idstock1)) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8")
		res, err := stmt.Exec()
		Perror(err)
		log.Println(res.RowsAffected())
	} else {
		log.Println("Will not create Tabel")
	}

	tableExist = true

}

//DBTestConnection is a simple method to test connection to DB
func DBTestConnection(dbtable string) {
	DbCreateConnectionString()

	//connString = "b06fa04f33f026:a480ce77@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_4499eae5e6a7bd4" //FUNKAR!!

	db, err := sql.Open(dbType, connString)
	Perror(err)

	var noRows int
	rows, err := db.Query("SELECT COUNT(*) FROM" + dbtable)
	log.Println(err)
	rows.Next()
	rows.Scan(&noRows)
	log.Println(noRows)
	db.Close()
}

//DbInsertSQL inserts a new stock into an sql server database
func DbInsertSQL(_stock models.StockDataSaveFormat, dbtable string) (string, error) {
	//if !tableExist {
	//DbTableCreate()
	//}

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	stmt, err := db.Prepare("INSERT " + dbtable + " SET symbol=?,created=?,buyprice=?,noshares=?,salesprice=?,name=?,lastprice=?, user=?")
	Perror(err)

	res, err := stmt.Exec(_stock.Symbol, _stock.Created, _stock.BuyPrice, _stock.NumberOfShares, _stock.SalesPrice, _stock.Name, _stock.LastTradePriceOnly)
	Perror(err)
	nid, err := res.LastInsertId()
	log.Println(nid)

	db.Close()

	log.Println("Insert i DB")

	return "Insert successfull", nil
}

//DBQuerySQL selects all stocks from the database
func DBQuerySQL(dbtable string) (Stocks, error) {
	if !tableExist {
		DbTableCreate(dbtable)
	}

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	var noRows int
	rows, err := db.Query("SELECT COUNT(*) FROM " + dbtable)
	Perror(err)
	for rows.Next() {
		err = rows.Scan(&noRows)
	}
	Perror(err)

	var newStocks = make(Stocks, noRows)

	rows, err = db.Query("SELECT * FROM " + dbtable)
	Perror(err)
	//log.Println(rows.Columns())
	var i int
	for rows.Next() {

		var tempStock models.StockDataSaveFormat
		//var id int
		err = rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly)
		Perror(err)
		//log.Println(tempStock.ID)
		newStocks[i] = tempStock
		//log.Println(newStocks[i].Symbol)
		i++
	}

	db.Close()
	log.Println("Hämtat data från DB")
	return newStocks, nil
}

//DBQuerySQLSingle selects one specific stocks from the database
func DBQuerySQLSingle(symbol string, dbtable string) (models.StockDataSaveFormat, error) {

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	stmt, err := db.Prepare("select * from " + dbtable + " where symbol = ?")
	defer stmt.Close()
	rows, err := stmt.Query(symbol)
	defer rows.Close()
	var tempStock models.StockDataSaveFormat
	//var id int
	for rows.Next() {
		err := rows.Scan(&tempStock.ID, &tempStock.Symbol, &tempStock.Created, &tempStock.BuyPrice, &tempStock.NumberOfShares, &tempStock.SalesPrice, &tempStock.Name, &tempStock.LastTradePriceOnly)
		Perror(err)
		//log.Println("Hämtar data från DB")
	}

	db.Close()
	return tempStock, nil
}

func DBQueryUserRelation(dbtable string, instocks Stocks) []models.User {

	DbCreateConnectionString()
	db, err := sql.Open(dbType, connString)
	Perror(err)

	//Stocks is a slice of StockDataSaveFormat
	type Users []models.User
	var userarray = make(Users, len(instocks))

	for i := 0; i < len(instocks); i++ {

		stmt, err := db.Prepare("select * from " + dbtable + " where id = (select userid from userrelation where stockid =?)")
		Perror(err)
		defer stmt.Close()
		rows, err := stmt.Query(instocks[i].ID)
		defer rows.Close()
		var tempUser models.User
		for rows.Next() {
			err := rows.Scan(&tempUser.User, &tempUser.Name, &tempUser.GoogleID)
			Perror(err)
			//log.Println(tempUser.Name)
			userarray[i] = tempUser

		}

	}

	return userarray

}

//DBDeletePost selects one specific stocks from the database
func DBDeletePost(dbtable string, symbol string) error {

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	stmt, err := db.Prepare("delete from " + dbtable + " where symbol=?")
	Perror(err)
	res, err := stmt.Exec(symbol)
	log.Println(res)

	db.Close()

	return err
}

//DBUpdateValue updates a selected field
func DBUpdateValue(dbtable string, symbol string, field string, value string) error {

	DbCreateConnectionString()

	db, err := sql.Open(dbType, connString)
	Perror(err)

	stmt, err := db.Prepare("UPDATE " + dbtable + " SET " + field + "=? WHERE symbol =?")
	Perror(err)

	res, err := stmt.Exec(value, symbol)
	Perror(err)
	log.Println(res)

	db.Close()

	return err

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

	dbUser = "b5fda560229b16"
	dbPassword = "10275c0d"
	dbServer = "tcp(us-cdbr-iron-east-04.cleardb.net:3306)"
	dbName = "heroku_949c9d7c3962055"
	dbType = "mysql"

}
