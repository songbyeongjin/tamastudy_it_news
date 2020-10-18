package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Connector struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Dbname   string `yaml:"db_name"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (d *Connector) GetConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, d.Dbname)
}

//Set Db information from yaml
func getDbConnector() (*Connector, error) {
	//temporary path for debug mode
	buf, err := ioutil.ReadFile("../../lib/external/db/db_info.yaml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var dbConnector Connector

	err = yaml.Unmarshal(buf, &dbConnector)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &dbConnector, nil
}



type Handler struct {
	Conn *gorm.DB
}

func NewDbHandler() (*Handler, error) {
	dbHandler := &Handler{}
	DbConnector, err := getDbConnector()

	if err != nil{
		log.Print(err)
		return nil, err
	}
	dbHandler.Conn, err = gorm.Open(DbConnector.Dialect, DbConnector.GetConnectString())

	if err != nil{
		log.Print(err)
		return nil, err
	}

	return dbHandler, nil
}

func NewTestDbHandler() *Handler {
	dbHandler := &Handler{}
	DbConnector, err := getDbConnector()


	if err != nil{
		log.Print(err)
		return nil
	}

	DbConnector.Dbname = DbConnector.Dbname + "_test"
	dbHandler.Conn, err = gorm.Open(DbConnector.Dialect, DbConnector.GetConnectString())

	if err != nil{
		log.Print(err)
		return nil
	}

	return dbHandler
}