package middlework

import(
	"log"
	"io"
	"strconv"
	"fmt"
	"sort"
	"math"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/shashank404error/shashankMongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pubnub/go"
)

var zoneInfo shashankMongo.ZoneInfo

type DistanceSorter []shashankMongo.DeliveryDetail //Interface to implement sorting of distance

func (a DistanceSorter) Len() int           { return len(a) }
func (a DistanceSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DistanceSorter) Less(i, j int) bool { return a[i].DistanceFromYou < a[j].DistanceFromYou }

func CreateZones(dBConnect *shashankMongo.ConnectToDataBase,collectionName string, userId string, config *shashankMongo.ProfileConfig) {
	for _,v := range config.ZoneID {

		load:=`{
			"name":"`+v+`",
			"businessUid": "`+userId+`",
			"deliveryInZone": "0",
			"picurl":"/static/images/userNotLogged.jpg",
			"longitude":"77.471906",
			"latitude":"23.160734"
			}`
		loadToJson:=ByteToJsonInterface(load)
		_=shashankMongo.InsertOne(dBConnect,collectionName,loadToJson)
	}
}

func UploadToExcel(file io.Reader,dBConnect *shashankMongo.ConnectToDataBase,collectionName string, userId string) (int64, int64, string){
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return 0,0,""
	}
	rows := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return 0,0,""
		}
	var arrOfDeliveryDetail []shashankMongo.DeliveryDetail	
	var count int
	for _, row := range rows {
		latitudeFloat, err := strconv.ParseFloat(row[3], 64); 
		if err != nil {
			fmt.Println(err) 
			return 0,0,""
		}
		longitudeFloat, err := strconv.ParseFloat(row[4], 64); 
		if err != nil {
			fmt.Println(err) 
			return 0,0,""
		}
		deliveryDetail:=shashankMongo.DeliveryDetail{
			CustomerName: row[0], 
			CustomerMob: row[1],
			Address: row[2],
			PicURL:"/static/images/userNotLogged.jpg",
			Latitude: latitudeFloat,
			Longitude: longitudeFloat,
			LongLat: row[4]+","+row[3],
		}
		arrOfDeliveryDetail = append(arrOfDeliveryDetail,deliveryDetail)
		count = count + 1
	}
	
	resultDocument:=shashankMongo.GetFieldByID(dBConnect,collectionName,userId)
	
	bsonBytes, _ := bson.Marshal(resultDocument)
	bson.Unmarshal(bsonBytes, &zoneInfo)

	oldCount, err := strconv.Atoi(zoneInfo.DeliveryInZone)
	if err != nil {
		fmt.Println(err)
		return 0,0,""
	}

	newCount:=count+oldCount
	countString := strconv.Itoa(newCount)

	businessAccount:=shashankMongo.FetchProfile(dBConnect, "businessAccounts", zoneInfo.BusinessUID)
	totalPending, err := strconv.Atoi(businessAccount.DeliveryPending)
	if err != nil {
		fmt.Println(err)
		return 0,0,""
	}
	newPending:= totalPending+count
	newPendingString := strconv.Itoa(newPending)
	_=shashankMongo.UpdateOneByID(dBConnect,"businessAccounts",zoneInfo.BusinessUID,"deliveryPending",newPendingString)	

	res1:=shashankMongo.UpdateDeliveryInfo(dBConnect,collectionName,userId,arrOfDeliveryDetail)
	res2:=shashankMongo.UpdateOneByID(dBConnect,collectionName,userId,"deliveryInZone", countString)
	return res1,res2,zoneInfo.BusinessUID
}

func SortZoneInfo(zoneInfo *shashankMongo.ZoneInfo ,userLong string, userLat string,token string) *shashankMongo.ZoneInfo {
	
	var arrOfDeliveryDetail []shashankMongo.DeliveryDetail
	n,err:=strconv.Atoi(zoneInfo.DeliveryInZone)
	if err!=nil {
		log.Fatal(err)
	}
	for i := 0; i < n; i++ {
		var result shashankMongo.MapBoxResp
		resp, err := http.Get("https://api.mapbox.com/directions/v5/mapbox/driving/"+userLong+","+userLat+";"+zoneInfo.DeliveryDetail[i].LongLat+"?access_token="+token)
		if err != nil{
			log.Fatal(err)
		}

		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data,&result)
		dis :=  math.Round(((result.Routes[0].Distance)/1000)*100)/100
		dur :=  math.Round(((result.Routes[0].Duration)/3600)*100)/100
		delivery:=shashankMongo.DeliveryDetail{
			Address: zoneInfo.DeliveryDetail[i].Address, 
			DistanceFromYou: dis,
			ETA: dur,
			PicURL: zoneInfo.DeliveryDetail[i].PicURL,
			LongLat: zoneInfo.DeliveryDetail[i].LongLat,
			CustomerName: zoneInfo.DeliveryDetail[i].CustomerName,
			CustomerMob: zoneInfo.DeliveryDetail[i].CustomerMob,
			Latitude: zoneInfo.DeliveryDetail[i].Latitude,
			Longitude: zoneInfo.DeliveryDetail[i].Longitude,
		}
		arrOfDeliveryDetail = append(arrOfDeliveryDetail,delivery)
	}

	sort.Sort(DistanceSorter(arrOfDeliveryDetail))
	zoneInfo.DeliveryDetail = arrOfDeliveryDetail
	return zoneInfo
}

func CreatePubnubChannel(pubnubCred shashankMongo.PubnubCredentials,channelName string) {
	config := pubnub.NewConfig()
    config.SubscribeKey = pubnubCred.SubscribeKey
    config.PublishKey = pubnubCred.PublishKey
	config.UUID = pubnubCred.UUIDPubnub
	
	pn := pubnub.NewPubNub(config)
   //   doneConnect := make(chan bool)
	donePublish := make(chan bool)
	
	msg := map[string]interface{}{
		"lat": "0",
		"lng":"0",
	}
	response, status, err := pn.Publish().Channel(channelName).Message(msg).Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response, status, err)

	<-donePublish
}


func ByteToJsonInterface(load string) map[string]interface{} {
	var loadArr = []byte(load)
    var loadToJson map[string]interface{}
    err := json.Unmarshal(loadArr, &loadToJson)
    if (err != nil) {
		log.Fatal(err)
	}
	return loadToJson
}