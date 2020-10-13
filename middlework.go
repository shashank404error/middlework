package middlework

import(
	"fmt"
	"github.com/shashank404error/shashankMongo"
)

func CreateZones(dBConnect *shashankMongo.ConnectToDataBase, userId string, config *shashankMongo.ProfileConfig) {
	for _,v := range config.ZoneID {
		fmt.Println(v)
	}
}