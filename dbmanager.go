package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const DATABASE = "nginxLogger2"
const C_USER = "user"
const C_LOG = "access_log"
const C_FILE = "file_position"

type FilePosition struct{
	ID bson.ObjectId
	Position int64
}

var user_collection, log_collection, file_collection *mgo.Collection
func initDB() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	user_collection = session.DB(DATABASE).C(C_USER)
	log_collection = session.DB(DATABASE).C(C_LOG)
	file_collection = session.DB(DATABASE).C(C_FILE)
}
func insertUser( username, password string) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user_collection.Insert(bson.M{
		"username" : username,
		"secret" : string(hashedPassword)})
}
func insertPosition(position int){
	update_query := bson.M{
		"$set": bson.M{"position": position}}
	_, err  := file_collection.Upsert(bson.M{}, update_query)
	if err != nil {
		panic(err)
	}
}
func getFilePointer() int64 {
	var record FilePosition
	err := file_collection.Find(nil).One(&record)
	if err != nil {
		panic(err)
	}
	return record.Position
}
func insertLog(rec Record){
	log_collection.Insert(&rec)
}