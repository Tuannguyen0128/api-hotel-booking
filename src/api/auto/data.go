package auto

import "ProjectPractice/src/api/models"

var teammembers = []models.TeamMember{
	{Fullname: "Tuan1", Email: "vt1@gmail.com", Password: "123456", MerchantCode: "MC101"},
	{Fullname: "Tuan2", Email: "vt2@gmail.com", Password: "123456", MerchantCode: "MC102"},
	{Fullname: "Tuan3", Email: "vt3@gmail.com", Password: "123456", MerchantCode: "MC101"},
	{Fullname: "Tuan4", Email: "vt4@gmail.com", Password: "123456", MerchantCode: "MC102"},
}
var merchantaccount = []models.MerchantAccount{
	{Description: "Merchant desciption 1.", MerchantCode: "MC101"},
	{Description: "Merchant desciption 2.", MerchantCode: "MC102"},
	{Description: "Merchant desciption 3.", MerchantCode: "MC103"},
	{Description: "Merchant desciption 4.", MerchantCode: "MC104"},
	{Description: "Merchant desciption 5.", MerchantCode: "MC105"},
}
