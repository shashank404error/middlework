package middlework

import(
	"log"
	"io"
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

	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    //fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)
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
	for _, row := range rows {
		fmt.Println(row[0])
		fmt.Println(row[1])
		fmt.Println(row[2])
		fmt.Println(row[3])
		fmt.Println(row[4])
		fmt.Println("\n")
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