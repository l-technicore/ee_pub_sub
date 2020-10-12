package model

import(
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "encoding/json"
  "asig/utils"
)

var db *gorm.DB
type offer struct {
	hotel hotel 				`json:"hotel"`
	room room 					`json:"room"`
	rate_plan rate_plan 		`json:"rate_plan"`
	cm_offer_id string 			`json:"cm_offer_id"`
    original_data string 		`json:"original_data"`
    capacity string 			`json:"capacity"`
    number json.Number 			`json:"number"`
    price json.Number 			`json:"price"`
    currency string 			`json:"currency"`
    check_in string 			`json:"check_in"`
    check_out string 			`json:"check_out"`
    fees []fees 				`json:"fees"`
}

type capacity struct{
	max_adults json.Number			`json:"max_adults"`
	extra_children json.Number		`json:"extra_children"`
}

type original_data struct {
	GuaranteePolicy GuaranteePolicy `json:"GuaranteePolicy"`
}

type GuaranteePolicy struct {
	Required bool 			`json:"Required"`
}

type fees struct{
	Type string				`json:"type"`
	description string		`json:"description"`
	included bool			`json:"included"`
	percent float64			`json:"percent"`
}

type hotel struct {
	hotel_id string				`json:"hotel_id"`
	name string					`json:"name"`
	country string				`json:"country"`
	address string				`json:"address"`
	latitude float64			`json:"latitude"`
	longitude float64			`json:"longitude"`
	telephone string			`json:"telephone"`
	amenities []string			`json:"amenities"`
	description string			`json:"description"`
	room_count float64			`json:"room_count"`
	currency string				`json:"currency"`
}

type room struct {
	hotel_id string				`json:"hotel_id"`
	room_id string				`json:"room_id"`
	description string			`json:"description"`
	name string					`json:"name"`
	max_adults float64			`json:"max_adults"`
	extra_children float64		`json:"extra_children"`
}

type rate_plan struct {
	hotel_id string								`json:"hotel_id"`
	rate_plan_id string							`json:"rate_plan_id"`
	cancellation_policy []cancellation_policy	`json:"cancellation_policy"`
	name string									`json:"name"`
	other_conditions []string					`json:"other_conditions"`
	meal_plan string							`json:"meal_plan"`
}

type Hotel struct {
	gorm.Model
	HotelId string				`gorm:"index"`
	name string								
	country string							
	address string							
	latitude float64						
	longitude float64						
	telephone string						
	amenities []string						
	description string						
	room_count float64						
	currency string				
	Room room					`gorm:"foreignKey:HotelId"`
	Rate_plan rate_plan			`gorm:"foreignKey:HotelId"`
}

type Room struct {
	gorm.Model
	HotelId string				`gorm:"index" json:"hotel_id"`
	room_id string				`gorm:"index"`
	description string						
	name string								
	max_adults float64						
	extra_children float64					
}

type Rate_plan struct {
	gorm.Model
	HotelId string								`gorm:"index" json:"hotel_id"`
	Rate_planId string							`gorm:"index"`
	cancellation_policy []cancellation_policy	`gorm:"foreignKey:Rate_planId"`
	name string												
	other_conditions []string					
	meal_plan string										
}

type cancellation_policy struct {
	Id string										
	Type string										
	expires_days_before string						
}

func InitializeDB(db_file string) (err error) {
	db, err = gorm.Open(sqlite.Open(db_file), &gorm.Config{PrepareStmt: true})
	if err != nil {
		utils.Log("failure", "995", "Unable to subscribe queue! Error:"+err.Error(), map[string]interface{}{"db_file": db_file})
		return
	}

	// Create Required Tables
	db.AutoMigrate(&Hotel{})
	// db.AutoMigrate(&room{})
	// db.AutoMigrate(&rate_plan{})

	// if !db.HasTable(&hotel{}){
	// 	db.CreateTable(&hotel{})
	// }
	// if !db.HasTable(&room{}){
	// 	db.CreateTable(&room{})
	// }
	// if !db.HasTable(&rate_plan{}){
	// 	db.CreateTable(&rate_plan{})
	// }
	return
}

// func Load(offers map[string]interface{}){
// 	hotel_amenities, _ := json.Marshal(offer["hotel"].(map[string]interface{})["amenities"])
// 	db.Create(&hotel{
// 				hotel_id:		offer["hotel"].(map[string]interface{})["hotel_id"].(string),
// 				name:			offer["hotel"].(map[string]interface{})["name"].(string),
// 				country:		offer["hotel"].(map[string]interface{})["country"].(string),
// 				address:		offer["hotel"].(map[string]interface{})["address"].(string),
// 				latitude:		offer["hotel"].(map[string]interface{})["latitude"].(float64),
// 				longitude:		offer["hotel"].(map[string]interface{})["longitude"].(float64),
// 				telephone:		offer["hotel"].(map[string]interface{})["telephone"].(string),
// 				amenities:		string(hotel_amenities),
// 				description:	offer["hotel"].(map[string]interface{})["description"].(string),
// 				room_count:		offer["hotel"].(map[string]interface{})["room_count"].(float64),
// 				currency:		offer["hotel"].(map[string]interface{})["currency"].(string),
// 			})
// 	db.Create(&room{
// 				hotel_id:		offer["room"].(map[string]interface{})["hotel_id"].(string),
// 				room_id:		offer["room"].(map[string]interface{})["room_id"].(string),
// 				description:	offer["room"].(map[string]interface{})["description"].(string),
// 				name:			offer["room"].(map[string]interface{})["name"].(string),
// 				max_adults:		offer["room"].(map[string]interface{})["capacity"].(map[string]interface{})["max_adults"].(float64),
// 				extra_children:	offer["room"].(map[string]interface{})["capacity"].(map[string]interface{})["extra_children"].(float64),
// 			})
// 	rate_plan_cancellation_policy, _ := json.Marshal(offer["rate_plan"].(map[string]interface{})["cancellation_policy"])
// 	rate_plan_other_conditions, _ := json.Marshal(offer["rate_plan"].(map[string]interface{})["other_conditions"])
// 	db.Create(&rate_plan{
// 				hotel_id:				offer["rate_plan"].(map[string]interface{})["hotel_id"].(string),
// 				rate_plan_id:			offer["rate_plan"].(map[string]interface{})["rate_plan_id"].(string),
// 				cancellation_policy:	string(rate_plan_cancellation_policy),
// 				name:					offer["rate_plan"].(map[string]interface{})["name"].(string),
// 				other_conditions:		string(rate_plan_other_conditions),
// 				meal_plan:				offer["rate_plan"].(map[string]interface{})["meal_plan"].(string),
// 			})
// }


func Load(body []byte) {
	var offers []offer
	if err := json.Unmarshal(body, &offers); err != nil {
		utils.Log("failure", "998", err.Error(), map[string]interface{}{"input": string(body)})
		return
	}
	var hotels []Hotel
	for _, o := range offers {
		var h Hotel
		temp, _ := json.Marshal(o.hotel)
		json.Unmarshal(temp, &h)
		// h = o.hotel
		temp, _ = json.Marshal(o.room)
		json.Unmarshal(temp, &h.Room)
		// h.Room = o.room
		temp, _ = json.Marshal(o.rate_plan)
		json.Unmarshal(temp, &h.Rate_plan)
		// h.Rate_plan = o.rate_plan
		hotels = append(hotels, h)
	}
	db.Create(hotels)
}