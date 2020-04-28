package graphql

import (
	"encoding/json"
	"fUber/types"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"pault.ag/go/haversine"
)

var cabType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Cab",
		Fields: graphql.Fields{
			"Cab_id": &graphql.Field{
				Type: graphql.Int,
			},
			"Cab_location": &graphql.Field{
				Type: graphql.String,
			},
			"Cab_lat": &graphql.Field{
				Type: graphql.Float,
			},
			"Cab_long": &graphql.Field{
				Type: graphql.Float,
			},
			"Cab_status": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{

			"cabs": &graphql.Field{
				Type:        graphql.NewList(cabType),
				Description: "Fetch cab with location",
				Args: graphql.FieldConfigArgument{
					"Cab_location": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var temp []types.Cab
					location, ok := p.Args["Cab_location"].(string)

					if ok {
						db, _ := gorm.Open("mysql", "root:root@tcp(mysqldb:3306)/fUber")
						Db := db.Table("cabs").Where("cab_location=?", location)
						Db.Where("cab_status=?", "available").Find(&temp)

						return temp, nil
					}
					return nil, nil
				},
			},
		},
	}, )
var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func Findcab( from string, lat float64, long float64, cabChan chan types.BestCabs) {
	query := fmt.Sprintf(`
{
cabs(Cab_location:"%v"){
Cab_id
Cab_location
Cab_lat
Cab_long
Cab_status
}
}
`, from)
	var bcab types.BestCabs

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		cabChan <- bcab

		return
	}

	res, _ := json.Marshal(r.Data)
	var v1 types.Cabs
	json.Unmarshal(res, &v1)
	//fmt.Println(v1)
	var dist = 5000.00
	Userlocation := haversine.Point{Lat: lat, Lon: long}

	for _, cab := range v1.Cabs {
		Cablocation := haversine.Point{Lat: cab.Cab_lat, Lon: cab.Cab_long}
		if float64(Userlocation.MetresTo(Cablocation)) <= dist {
			bcab.Distance = float64(Userlocation.MetresTo(Cablocation))
			bcab.Cab_id = cab.Cab_id
				//fmt.Println(bcab)
		}
	}
		cabChan <- bcab
	return
}