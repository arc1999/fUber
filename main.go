package main

import (
	"errors"
	"fUber/Auth"
	"fUber/Request"
	"fUber/Response"
	"fUber/errs"
	"fUber/graphql"
	"fUber/types"
	"fUber/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)
type Mysql struct {
	Db *gorm.DB
}

func main() {
	dba, err := gorm.Open("mysql","root:root@tcp(mysqldb:3306)/")
	dba.Exec("CREATE DATABASE IF NOT EXISTS"+" fUber")
	dba.Close()

	db, err := gorm.Open("mysql", "root:root@tcp(mysqldb:3306)/fUber")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	if (!db.HasTable(&types.Booking{})) {
		db.AutoMigrate(&types.Booking{})
	}
	if (!db.HasTable(&types.User{})) {
		db.AutoMigrate(&types.User{})
	}
	if (!db.HasTable(&types.Cab{})) {
		db.AutoMigrate(&types.Cab{})
	}
	set := &Mysql{db}
	routes(set)
	defer db.Close()
}

func routes(set * Mysql){
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/login", func(r chi.Router) {
		r.Post("/",Auth.TokenHandler)
		r.Post("/signup",Auth.Signup)
	})
	r.Route("/booking", func(r chi.Router) {
		r.Use(Auth.AuthMiddleware)
		r.Post("/", set.Booking)
		r.Get("/", set.ListAllBookings)
		r.Post("/cabs",set.cabs)
	})
	log.Fatal(http.ListenAndServe(":8080", r))

}

func (Db *Mysql)ListAllBookings(writer http.ResponseWriter ,request *http.Request){
	token := request.Context().Value("user").(*jwt.Token)
	a:=token.Claims.(jwt.MapClaims)
	id:=int(a["user"].(float64))
	var temps types.Bookings
	Dba:=Db.Db.Table("bookings").Where("u_id=?",id).Find(&temps)
	if Dba.RowsAffected == 0{
		err:=errors.New("Bookings not found")
		_ = render.Render(writer, request, errs.ErrRender(err))
		return
	}
	render.Render(writer,request,Response.ListBookings(&temps))
}

func (Db *Mysql) Booking(writer http.ResponseWriter, request *http.Request) {
	token := request.Context().Value("user").(*jwt.Token)
	a := token.Claims.(jwt.MapClaims)
	id := int(a["user"].(float64))
	var booking1 Request.CreateBookingRequest
	err := render.Bind(request, &booking1)
	if err != nil {
		_ = render.Render(writer, request, errs.ErrRender(err))
	}
	City := [3][3]string{}
	City[0][0] = "Gultekdi" //18.494726  73.8669
	City[0][1] = "BoradeNagar" //18.4847 73.8926
	City[0][2] = "SayadNagar" //18.4865 , 73.9294
	City[1][0] = "Bibwewadi" //18.4756 , 73.8622
	City[1][1] = "Kondhwa" //18.4695 , 73.8890
	City[1][2] = "Hadapsar"  //18.4751 , 73.9163
	City[2][0] = "MadhalaVada" //18.4550 , 73.8856
	City[2][1] = "Pisoli" //18.4492 , 73.9078
	City[2][2] = "PandhariNagar" //18.5536 , 73.9407
	var (
		i, j  int
		spawnedRoutines int
		threshold_dist = 5000.0
		all_cabs []types.BestCabs
		cab types.BestCabs
	)
	cabChannel := make(chan types.BestCabs)


	// check for the current block in which the user is located



	go graphql.Findcab( booking1.From, booking1.U_lat, booking1.U_long, cabChannel)
	nearestCab := <- cabChannel // get the matching channel

	fmt.Println(nearestCab)

	//check for nearby locality

	if nearestCab.Cab_id == 0 {
		// spawn the routines and check in all blocks except the current one
		for i = 0; i < 3; i++ {
			for j = 0; j < 3; j++ {
				if City[i][j] == booking1.From {
					goto Exit
				}
			}
		}
	Exit:
		if(i>=3 && j>=3){
			fmt.Fprint(writer,"Not Serving in your Area")
			return
		}
		if i+1 < 3 {
			spawnedRoutines++

			go graphql.Findcab( City[i+1][j], booking1.U_lat, booking1.U_long, cabChannel)
			if j+1 < 3 {
				spawnedRoutines++

				go graphql.Findcab(City[i+1][j+1], booking1.U_lat, booking1.U_long, cabChannel)
			}
			if j-1 >= 0 {
				spawnedRoutines++

				go graphql.Findcab( City[i+1][j-1], booking1.U_lat, booking1.U_long, cabChannel)
			}
		}
		if i-1 >= 0 {
			spawnedRoutines++

			go graphql.Findcab( City[i-1][j], booking1.U_lat, booking1.U_long, cabChannel)
			if j+1 < 3 {
				spawnedRoutines++

				go graphql.Findcab( City[i-1][j+1], booking1.U_lat, booking1.U_long, cabChannel)
				spawnedRoutines++

				go graphql.Findcab( City[i][j+1], booking1.U_lat, booking1.U_long, cabChannel)
			}
			if j-1 >= 0 {
				spawnedRoutines++

				go graphql.Findcab( City[i-1][j-1], booking1.U_lat, booking1.U_long, cabChannel)
				spawnedRoutines++

				go graphql.Findcab( City[i][j-1], booking1.U_lat, booking1.U_long, cabChannel)
			}
		}

		for i := 0; i < spawnedRoutines; i++ {

			select {
			case cab = <-cabChannel:
				if(cab.Cab_id>0){
					all_cabs=append(all_cabs,cab)
				}

				// if any cab id is returned from any func then use that to book the cab

			}
		}
		close(cabChannel)
		//wg.Wait()
		fmt.Println(all_cabs)
		for _,cab := range all_cabs {
			if (cab.Distance< threshold_dist) {
				threshold_dist = cab.Distance
				nearestCab = cab
			}
		}
		fmt.Println(nearestCab)
		if nearestCab.Cab_id == 0 {
			fmt.Fprint(writer, "No Cabs in your Area")
			return
		}
	}
	Db.Db.Table("cabs").Where("cab_id=?", nearestCab.Cab_id).Update(types.Cab{Cab_status: "unavailable"})

	booking1.Cabid = nearestCab.Cab_id
	booking1.Amount =  (nearestCab.Distance)*0.008
	booking1.U_id = id
	Db.Db.Create(&booking1.Booking)
	render.Render(writer, request, Response.ListBooking(&booking1))
}
func (Db *Mysql)cabs(writer http.ResponseWriter ,request *http.Request){
	for _,cab := range utils.Init_cabs {
		Dba := Db.Db.Create(&cab)

		if Dba.RowsAffected == 0 {
			err := errors.New("Cab  not initialized")
			_ = render.Render(writer, request, errs.ErrRender(err))
			return
		}
	}
}
