package db

var DB *Databases

type Databases struct {
	MySQL *MysqlDB
}

func init() {
	DB = new(Databases)
	sql, err := InitMySQL()
	if err != nil {
		panic(err)
	}
	DB.MySQL = sql
}
