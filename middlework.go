package middlework

import(
	"log"
	"io"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/shashank404error/shashankMongo"
	"github.com/360EntSecGroup-Skylar/excelize"
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

func UploadToExcel(file io.Reader,dBConnect *shashankMongo.ConnectToDataBase,collectionName string, userId string) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return
		}
	var arrOfDeliveryDetail []shashankMongo.DeliveryDetail	
	for _, row := range rows {
		latitudeFloat, err := strconv.ParseFloat(row[3], 64); 
		if err != nil {
			fmt.Println(err) 
			return
		}
		longitudeFloat, err := strconv.ParseFloat(row[4], 64); 
		if err != nil {
			fmt.Println(err) 
			return
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
	}
	res:=shashankMongo.UpdateDeliveryInfo(dBConnect,collectionName,userId,arrOfDeliveryDetail)
	fmt.Println(res)
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