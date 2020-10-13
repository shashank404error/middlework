package middlework

import(
	"fmt"
)

type ConnectToDataBase struct {
	CustomApplyURI string 
	DatabaseName string 
	CollectionName string 
}

type ProfileConfig struct{
	Zone int64 `bson: "zone" json: "zone"`
	MessagePlan int64 `bson: "messageplan" json: "messageplan"`
	Tracking bool `bson: "tracking" json: "tracking"`
	ZoneID []string `bson: "zoneid" json: "zoneid"`
}

var dBConnect *ConnectToDataBase
var config *ProfileConfig

func CreateZones(dBConnect *ConnectToDataBase, userId string, config *ProfileConfig) {
	for _,v := range config.ZoneID {
		fmt.Println(v)
	}
}