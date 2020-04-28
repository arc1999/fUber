package types

type Booking struct {
	Bookingid int `gorm:"primary_key; unique_index; not null" json:"bookingid"`
	From string `gorm:"not null" json:"from"`
	To string `gorm:"not null" json:"to"`
	Cabid int `gorm:"not null" json:"cabid"`
	Amount float64 `gorm:"not null" json:"amount"`
	U_id int `gorm:"not null" json:"u_id"`
	U_lat float64 `gorm:"not null" json:"u_lat"`
	U_long float64 `gorm:"not null" json:"u_long"`
}
type Bookings []Booking
type Cab struct {
	Cab_id int `gorm:"primary_key; unique_index; not null" json:"cab_id"`
	Cab_location string `gorm:"not null" json:"cab_location"`
	Cab_lat float64 `gorm:"not null" json:"cab_lat"`
	Cab_long float64 `gorm:"not null" json:"cab_long"`
	Cab_status string `gorm:"not null" json:"cab_status"`
}
type Cabs struct {
	Cabs []Cab `json:"cabs"`
}
type User struct{
	U_id   int    `gorm:"primary_key; unique_index; not null" json:"u_id"`
	U_name string `gorm:"primary_key; not null" json:"u_name"`
	U_pass string `gorm:"not null" json:"u_pass"`
	Amount float64 `gorm:"not null" json:"amount"`
	Email_id string `gorm:"not null" json:"email_id"`

}
type BestCabs struct {
	Cab_id int `json:"cab_id"`
	Distance float64 `json:"distance"`
}
