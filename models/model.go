package models

import (
	"time"

	"github.com/lib/pq"
)

type ChatAssistant struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	UserID    string    `json:"userid"`
	CreatedAt time.Time `json:"createdat"`
}

type Vendor struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	BusinessName    string         `json:"business_name"`
	BusinessAddress string         `json:"business_address"`
	PhoneNumber     string         `json:"phone_number"`
	UniqueID        string         `json:"unique_id"`
	Token           string         `json:"token"`
	District        pq.StringArray `json:"district"`
	Panchayat       pq.StringArray `json:"panchayat"`
	Services        pq.StringArray `json:"services"`
	ServiceCategory pq.StringArray `json:"service_category"`
	Experience      string         `json:"experience"`
	Reference       string         `json:"reference"`
	Premium         string         `json:"premium"`
	ProfileImage    string         `json:"profile_image"`
	PlanPurchase    string         `json:"plan_purchase"`
	Image           string         `json:"image"`
	ImageView       string         `json:"image_view"`
	Favorites       string         `json:"favorites"`
	CreatedAt       string         `json:"created_at"`
	Latitude        string         `json:"latitude"`
	Longitude       string         `json:"longitude"`
}

type Customer struct {
	ID                 uint   `json:"id"`
	Auto_taxi_goods_id string `json:"auto_taxi_goods_id"`
	Name               string `json:"name"`
	Phone_number       string `json:"phone_number"`
	Userid             string `json:"userid"`
	Token              string `json:"token"`
	District           string `json:"district"`
	Panchayat          string `json:"panchayat"`
	Profile_image      string `json:"profile_image"`
	Join_date          string `json:"join_date"`
	Is_block           bool   `json:"is_block"`
}

type Banner struct {
	ID               uint      `json:"id"`
	Image            string    `json:"image"`
	District         string    `json:"district"`
	Panchayat        string    `json:"panchayat"`
	Services         string    `json:"services"`
	Service_category string    `json:"service_category"`
	Created_at       time.Time `json:"createdat"`
}

type Ratings struct {
	ID         uint      `json:"id"`
	Created_at time.Time `json:"created_at"`
	Rating     string    `json:"rating"`
	Vendor_id  string    `json:"vendor_id"`
	Remark     string    `json:"remark"`
}

type AdRequest struct {
	ID           uint   `json:"id"`
	Ad_category  string `json:"ad_category"`
	Phone_number string `json:"phone_number"`
	Unique_id    string `json:"unique_id"`
	Category     string `json:"category"`
	District     string `json:"district"`
	Start_date   string `json:"start_date"`
	Request_date string `json:"request_date"`
}

type AutoTaxiGoods struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Phone_number  string    `json:"phone_number"`
	Unique_id     string    `json:"unique_id"`
	Service_type  string    `json:"service_type"`
	Experience    string    `json:"experience"`
	District      string    `json:"district"`
	Panchayat     string    `json:"panchayat"`
	Profile_image string    `json:"profile_image"`
	Image         string    `json:"image"`
	Image_views   string    `json:"image_views"`
	Rating        int64     `json:"rating"`
	Created_at    time.Time `json:"created_at"`
}

type Orders struct {
	ID                uint      `json:"id"`
	Order_id          string    `json:"order_id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Category          string    `json:"category"`
	Created_at        time.Time `json:"created_at"`
	Description       string    `json:"description"`
	Assigned_vendors  string    `json:"assigned_vendors"`
	End_date          time.Time `json:"end_date"`
	Image             string    `json:"image"`
	Service           string    `json:"service"`
	Sqf               string    `json:"sqf"`
	Start_date        time.Time `json:"start_date"`
	User_id           string    `json:"user_id"`
	Phone_number      string    `json:"phone_number"`
	Population        string    `json:"population"`
	Function_type     string    `json:"function_type"`
	Ac_available      string    `json:"ac_available"`
	Food_available    string    `json:"food_available"`
	Event_date        string    `json:"event_date"`
	Program_type      string    `json:"program_type"`
	Travel_experience string    `json:"travel_experience"`
	Origin            string    `json:"origin"`
	Destination       string    `json:"destination"`
	Vehicle_type      string    `json:"vehicle_type"`
	Dining            string    `json:"dining"`
	Accommodation     string    `json:"accommodation"`
	Property_type     string    `json:"property_type"`
	Budget            string    `json:"budget"`
	Property_location string    `json:"property_location"`
	Area              string    `json:"Area"`
	Quantity          string    `json:"quantity"`
	Tool_type         string    `json:"tool_type"`
	Seat_capacity     string    `json:"seat_capacity"`
	Rent_period       string    `json:"rent_period"`
	Service_type      string    `json:"service_type"`
	Type              string    `json:"type"`
	Product_type      string    `json:"product_type"`
	Course_type       string    `json:"course_type"`
	Called_vendors    string    `json:"called_vendors"`
	Building_type     string    `json:"building_type"`
	IsRating_Enable   bool      `json:"isRating_enable"`
	Property_address  string    `json:"property_address"`
}

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	CreatedAt    string `json:"created_at"`
	Is_vendor    bool   `json:"is_vendor"`
}

type Log struct {
	ID              uint   `json:"id"`
	CreatedDate     string `json:"created_date"`
	Log_Number      string `json:"log_number"`
	Log_Description string `json:"log_description"`
	Function        string `json:"function"`
	Estimate_Time   string `json:"estimate_time"`
}

type App struct {
	ID               uint   `json:"id"`
	Log_Number       string `json:"log_number"`
	Log_Description  string `json:"log_description"`
	Function         string `json:"function"`
	CreatedDate      string `json:"created_date"`
	Userid           string `json:"userid"`
	Device           string `json:"device"`
	Platform         string `json:"platform"`
	Estimate_Time    string `json:"estimate_time"`
	Exceptional_Case string `json:"exception_case"`
}
