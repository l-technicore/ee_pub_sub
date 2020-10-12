1. Extract zip in /tmp/ directory.

2. Run on Bash:
	export GOPATH=/tmp/ee_pub_sub
	cd /tmp/ee_pub_sub/src/asig/pub
	go get -v ./...
	sudo apt install rabbitmq-server

3. Publisher:
	main() => /tmp/ee_pub_sub/src/asig/pub/pub.go
	To run Publisher Execute the following in Bash:
		cd /tmp/ee_pub_sub/src/asig/pub
		go run pub.go
	Local Packages Used:
		"asig/publisher"
		"asig/pub/input"
		"asig/utils"

4. Subscriber:
	main() => /tmp/ee_pub_sub/src/asig/sub/sub.go
	To run Subscriber Execute the following in Bash:
		cd /tmp/ee_pub_sub/src/asig/sub
		go run sub.go
	Local Packages Used:
		"asig/subscriber"
		"asig/utils"
		"asig/model"


# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # 
#	 May be due to my forst time using gorm in Golang,
#	 I was unable to get "gorm.io/gorm" package to do the most basic table creation operation
#	 So, just wrote what I though should have worked to rech the goal on using sqlite DB.
#
#	 i.e. I tried using following structure:
#		type Hotel struct {
#			gorm.Model
#			hotel_id string				`gorm:"index"`
#			name string								
#			country string							
#			address string							
#			latitude float64						
#			longitude float64						
#			telephone string						
#			amenities []string						
#			description string						
#			room_count float64						
#			currency string				
#		}
#	 & using following command to try creating a corresponding table on sqlite:
#    	db.AutoMigrate(&Hotel{})
#	 But the table being created lacked all fields except the default ID, CreatedAt, UpdatedAt & DeletedAt fields,
#	 unlike what is being mentioned in their doc.
#    Refrence: https://gorm.io/docs/models.html
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #


5. Local Packages and their _test.go files:
	package sub: "asig/subscriber/subscriber.go" & "asig/subscriber/subscriber_test.go"
	package pub: "asig/publisher/publisher.go" & "asig/publisher/publisher_test.go"
	package utils: "asig/utils/logger.go" & "asig/utils/logger_test.go"
	package model: "asig/model/db.go"

6. To Test the packages Run "go test" on Bash in following directories:
	/tmp/ee_pub_sub/src/asig/utils/
	/tmp/ee_pub_sub/src/asig/subscriber/
	/tmp/ee_pub_sub/src/asig/publisher/