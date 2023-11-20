package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deeksharma/bookings/internal/config"
	"github.com/Deeksharma/bookings/internal/driver"
	"github.com/Deeksharma/bookings/internal/forms"
	"github.com/Deeksharma/bookings/internal/helpers"
	"github.com/Deeksharma/bookings/internal/models"
	"github.com/Deeksharma/bookings/internal/render"
	"github.com/Deeksharma/bookings/internal/repository"
	"github.com/Deeksharma/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Repository pattern, to swap variables between repos with minimal code changes
// share database connections using repository pattern
// all the receivers will have the same data

// Repo the repository used by handlers
var Repo *Repository

var layout = "02-01-2006"

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new rep
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestingRepo creates a new rep
func NewTestingRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// ****We can also use other templating engines other than .page.tmpl

// Home is the home page handler
func (m *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	m.DB.AllUsers()
	// saving session at the beginning
	remoteIP := request.RemoteAddr
	m.App.SessionManager.Put(request.Context(), "remote_ip", remoteIP)
	render.Template(writer, request, "home.page.tmpl", &models.TemplateData{})
}

// About is the home page handler
func (m *Repository) About(writer http.ResponseWriter, request *http.Request) {
	// perform some logic which gives data

	remoteIP := m.App.SessionManager.GetString(request.Context(), "remote_ip") // will have empty value if not present in session
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again!!" + remoteIP
	stringMap["remote_ip"] = remoteIP
	//m.App.SessionManager.Get() - access to session variable now

	// send data to the template
	render.Template(writer, request, "about.page.tmpl", &models.TemplateData{
		StringData: stringMap,
	})
}

// Generals renders the room page
func (m *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "generals-quarter.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "majors-suite.page.tmpl", &models.TemplateData{})
}

// MakeReservation renders the reservation page
func (m *Repository) MakeReservation(writer http.ResponseWriter, request *http.Request) {
	res, ok := m.App.SessionManager.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.SessionManager.Put(request.Context(), "error", "can't get reservation from session")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(writer, errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomById(context.Background(), res.RoomId)
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		// helpers.ServerError(writer, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.SessionManager.Put(request.Context(), "reservation", res)

	sd := res.StartDate.Format(layout)
	ed := res.EndDate.Format(layout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(writer, request, "make-reservation.page.tmpl", &models.TemplateData{
		Form:       forms.NewForm(nil), // we'll have the form object first time the object was rendered
		Data:       data,
		StringData: stringMap,
	})
}

// PostReservation handles posting of a reservation form
func (m *Repository) PostReservation(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := m.App.SessionManager.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.SessionManager.Put(request.Context(), "error", "can't get reservation from session")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(writer, errors.New("cant get reservation from session"))
		return
	}
	// if we have errors in the for validation then we'll render the form again but this time we'll show errors
	err := request.ParseForm()
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(writer, err)
		return
	}

	reservation.FirstName = request.Form.Get("first_name")
	reservation.LastName = request.Form.Get("last_name")
	reservation.Email = request.Form.Get("email")
	reservation.Phone = request.Form.Get("phone")

	form := forms.NewForm(request.PostForm) // request.PostForm is of the form url.Values

	form.Required("first_name", "last_name", "email")
	form.MinLength("phone", 10)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(writer, form.Errors.Get("first_name")+form.Errors.Get("last_name")+form.Errors.Get("email"), http.StatusSeeOther)
		render.Template(writer, request, "make-reservation.page.tmpl", &models.TemplateData{ // render the template again with the previous form data
			Form: form, // we'll have the form object first time the object was rendered
			Data: data,
		})
		return
	}
	ctx := context.Background()

	reservationId, err := m.DB.InsertReservation(ctx, reservation)
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	// now we need to put data in room restrictions that this room is now not available
	restriction := models.RoomRestriction{
		RoomId:        reservation.RoomId,
		ReservationId: reservationId,
		RestrictionId: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	err = m.DB.InsertRoomRestriction(ctx, restriction)
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(writer, err)
		return
	}

	// sending email to guest
	htmlMsg := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br><br>
		Dear %s %s,<br><br>
			This is to confirm your reservation from %s to %s for room %s.<br><br>
		Thank you,<br>
		XYZ
	`, reservation.FirstName, reservation.LastName, reservation.StartDate.Format(layout),
		reservation.EndDate.Format(layout), reservation.Room.RoomName)
	msg := models.MailData{
		To:       []string{reservation.Email},
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	m.App.MailChannel <- msg

	// sending email to owner
	htmlMsg = fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br><br>
		Dear Owner,<br><br>
			This is to confirm that your property's room %s has been booked from %s to %s by %s %s.<br><br>
		Thank you,<br>
		XYZ
	`, reservation.Room.RoomName, reservation.StartDate.Format(layout),
		reservation.EndDate.Format(layout), reservation.FirstName, reservation.LastName)

	msg = models.MailData{
		To:       []string{"me@here.com"},
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}

	m.App.MailChannel <- msg

	// we can use session to send variable reservation in the next page and retrieve it there
	m.App.SessionManager.Put(request.Context(), "reservation", reservation) // we are putting the session in the session variable, we should remove it later

	// then we'll redirect the page
	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)
}

// SearchAvailability renders the reservation page
func (m *Repository) SearchAvailability(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "reservation.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability renders the reservation page
func (m *Repository) PostSearchAvailability(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	start := request.Form.Get("start-date") //you always get a string from the form get data, need to typecast it later
	end := request.Form.Get("end-date")

	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		//helpers.ServerError(writer, err)
		return
	}
	endDate, _ := time.Parse(layout, end)
	if err != nil {
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		helpers.ServerError(writer, err)
		return
	}

	ctx := context.Background()
	rooms, err := m.DB.SearchAvailabilityForAllRoomsByDates(ctx, startDate, endDate)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	if len(rooms) == 0 {
		m.App.SessionManager.Put(request.Context(), "error", "No Availability...")
		http.Redirect(writer, request, "/reservation", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	// store the reservation model in session
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    0,
	}

	m.App.SessionManager.Put(request.Context(), "reservation", res)

	render.Template(writer, request, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
	//writer.Write([]byte("start date is" + start + "end date is" + end)) // this didnt work because it nosurf will ignore all the post requests that does not have CSRF token
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and sends a JSON response
func (m *Repository) AvailabilityJSON(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal Server Error",
		}
		out, _ := json.MarshalIndent(resp, "", "		")
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(out)
		return
	}

	sd := request.Form.Get("start")
	ed := request.Form.Get("end")

	roomIdString := request.Form.Get("room_id")
	roomId, err := strconv.Atoi(roomIdString)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	enDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	available, err := m.DB.SearchAvailabilityForRoomByDates(context.Background(), startDate, enDate, roomId)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to DB",
		}
		out, _ := json.MarshalIndent(resp, "", "		")
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(out)
		return
	}

	resp := jsonResponse{ // we are manually making a json response, so it'll always be right
		OK:        available,
		StartDate: sd,
		EndDate:   ed,
		RoomId:    roomIdString,
	}

	if resp.OK {
		resp.Message = "Available!!!"
	} else {
		resp.Message = "Not Available!!! :-("
	}

	out, _ := json.MarshalIndent(resp, "", "		")
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(out)
}

// Contact renders the reservation page
func (m *Repository) Contact(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary renders the reservation summary page
func (m *Repository) ReservationSummary(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := m.App.SessionManager.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		// search engines should not index this page and this page should be redirected to url not found or something
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.SessionManager.Put(request.Context(), "error", "Can't get reservation from session")
		log.Println("cannot get item from session")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.SessionManager.Remove(request.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format(layout)
	ed := reservation.EndDate.Format(layout)

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(writer, request, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:       data,
		StringData: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	res, ok := m.App.SessionManager.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(writer, errors.New("cannot get reservation from session"))
		return
	}

	res.RoomId = roomID
	m.App.SessionManager.Put(request.Context(), "reservation", res)
	http.Redirect(writer, request, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes url params, build a session variable and takes user to make-reservation
func (m Repository) BookRoom(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	sd := request.URL.Query().Get("s")
	startDate, err := time.Parse(layout, sd)
	ed := request.URL.Query().Get("e")
	endDate, err := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomById(context.Background(), roomID)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	res.RoomId = roomID
	res.Room.RoomName = room.RoomName
	res.EndDate = endDate
	res.StartDate = startDate

	m.App.SessionManager.Put(request.Context(), "reservation", res)
	log.Println(roomID, startDate, endDate)
	http.Redirect(writer, request, "/make-reservation", http.StatusTemporaryRedirect)

}

// ShowLogin is the get handler for user login
func (m *Repository) ShowLogin(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "login.page.tmpl", &models.TemplateData{
		Form: forms.NewForm(nil),
	})
}

// PostShowLogin is the post handler for user login
func (m *Repository) PostShowLogin(writer http.ResponseWriter, request *http.Request) {
	_ = m.App.SessionManager.RenewToken(request.Context()) // prevents session fixation attack - a token is fixed for a session so its always a good practice to renew token in case of login and logout as well

	err := request.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
		return
	}
	email := request.Form.Get("email")
	password := request.Form.Get("password")

	form := forms.NewForm(request.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		// TODO: take to usr/login again
		render.Template(writer, request, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(context.Background(), email, password)

	if err != nil {
		log.Println(err)
		m.App.SessionManager.Put(request.Context(), "error", err.Error())
		http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
		return
	}

	// now how to make the user login
	// we'll put the user id in the session
	m.App.SessionManager.Put(request.Context(), "user_id", id)
	m.App.SessionManager.Put(request.Context(), "flash", "Logged in successfully!")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

// Logout logouts a user
func (m *Repository) Logout(writer http.ResponseWriter, request *http.Request) {
	_ = m.App.SessionManager.Destroy(request.Context())
	_ = m.App.SessionManager.RenewToken(request.Context())
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

// AdminDashboard renders the admin dashboard page
func (m *Repository) AdminDashboard(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservation shows all new reservations that are not processed yet
func (m *Repository) AdminNewReservation(writer http.ResponseWriter, request *http.Request) {
	reservations, err := m.DB.AllNewReservations(context.Background())
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservations
	render.Template(writer, request, "admin-new-reservation.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservation shows all reservations
func (m *Repository) AdminAllReservation(writer http.ResponseWriter, request *http.Request) {
	reservations, err := m.DB.AllReservations(context.Background())
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservations
	render.Template(writer, request, "admin-all-reservation.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminReservationsCalender displays the reservation calender
func (m *Repository) AdminReservationsCalender(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "admin-reservations-calender.page.tmpl", &models.TemplateData{})
}

// display the reservation as form so that owner can make changes/ mark a reservation as processed

// AdminShowReservation shows the reservation in admin tool
func (m *Repository) AdminShowReservation(writer http.ResponseWriter, request *http.Request) {
	exploded := strings.Split(request.RequestURI, "/")
	source := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = source
	reservationId, err := strconv.Atoi(exploded[4])

	if err != nil {
		log.Println(err)
		helpers.ServerError(writer, err)
		return
	}

	log.Println("we are getting request from:", source, "reservation page, reservationId:", reservationId)

	reservation, err := m.DB.GetReservationById(context.Background(), reservationId)
	if err != nil {
		log.Println(err)
		helpers.ServerError(writer, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(writer, request, "admin-all-reservation-show.page.tmpl", &models.TemplateData{
		StringData: stringMap,
		Data:       data,
		Form:       forms.NewForm(nil),
	})
}

// Divide is the divide page handler
func (m *Repository) Divide(writer http.ResponseWriter, request *http.Request) {
	res, err := divide(2.0, 0.0)
	if err != nil {
		_, _ = fmt.Fprintf(writer, fmt.Sprint("Deeksha loves Yashwardhan and Yashwardhan loves Deeksha!!!!", err))
		return
	}
	_, _ = fmt.Fprintf(writer, fmt.Sprint("Deeksha loves Yashwardhan and Yashwardhan loves Deeksha!!!!", res))
}

// add adds two integer values
func add(x, y int) int {
	return x + y
}

func divide(x, y float32) (float32, error) {
	if y == 0 {
		return 0.0, fmt.Errorf("error divide by zero")
	}
	return x / y, nil
}
