package middlework

import(
	"log"
	"encoding/json"
	"github.com/shashank404error/shashankMongo"
)

func CreateZones(dBConnect *shashankMongo.ConnectToDataBase,collectionName string, userId string, config *shashankMongo.ProfileConfig) {
	for _,v := range config.ZoneID {

		load:=`{
			"name":"`+v+`",
			"businessUid": "`+userId+`",
			"deliveryInZone": "0"
			}`
		loadToJson:=byteToJsonInterface(load)
		_=shashankMongo.InsertOne(dBConnect,collectionName,loadToJson)
	}
}

func byteToJsonInterface(load string) map[string]interface{} {
	var loadArr = []byte(load)
    var loadToJson map[string]interface{}
    err := json.Unmarshal(loadArr, &loadToJson)
    if (err != nil) {
		log.Fatal(err)
	}
	return loadToJson
}