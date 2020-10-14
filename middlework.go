package middlework

import(
	"log"
	"io"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/shashank404error/shashankMongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/360EntSecGroup-Skylar/excelize"
)

var zoneInfo shashankMongo.ZoneInfo

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

func UploadToExcel(file io.Reader,dBConnect *shashankMongo.ConnectToDataBase,collectionName string, userId string) (int64, int64){
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return 0,0
	}
	rows := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return 0,0
		}
	var arrOfDeliveryDetail []shashankMongo.DeliveryDetail	
	var count int
	for _, row := range rows {
		latitudeFloat, err := strconv.ParseFloat(row[3], 64); 
		if err != nil {
			fmt.Println(err) 
			return 0,0
		}
		longitudeFloat, err := strconv.ParseFloat(row[4], 64); 
		if err != nil {
			fmt.Println(err) 
			return 0,0
		}
		deliveryDetail:=shashankMongo.DeliveryDetail{
			CustomerName: row[0], 
			CustomerMob: row[1],
			Address: row[2],
			Latitude: latitudeFloat,
			Longitude: longitudeFloat,
			LongLat: row[4]+","+row[3],
		}
		arrOfDeliveryDetail = append(arrOfDeliveryDetail,deliveryDetail)
		count = count + 1
	}
	countString := strconv.Itoa(count)
	resultDocument:=shashankMongo.GetFieldByID(dBConnect,collectionName,userId)
	
	bsonBytes, _ := bson.Marshal(resultDocument)
	bson.Unmarshal(bsonBytes, &zoneInfo)
	fmt.Println(zoneInfo)
	
	res1:=shashankMongo.UpdateDeliveryInfo(dBConnect,collectionName,userId,arrOfDeliveryDetail)
	res2:=shashankMongo.UpdateOneByID(dBConnect,collectionName,userId,"deliveryInZone", countString)
	return res1,res2
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