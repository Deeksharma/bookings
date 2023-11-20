package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Deeksharma/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TODO: Revisit the test again

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarter", "/generals-quarter", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK}, // removed post data from everywhere because the post functions is using sessions, so we'll send the data there only
	//{"reservation-summary", "/reservation-summary", "GET", http.StatusOK},

	//{"make-reservation", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value:"Deeksha"},
	//	{key: "last_name", value:"Sharma"},
	//	{key: "email", value:"deeksha.sharma@zomato.com"},
	//	{key: "phone", value:"07766909220"},
	//}, http.StatusOK},
	//{"reservation", "/reservation", "POST", []postData{
	//	{key: "start-date", value: "12-07-2022"},
	//	{key: "end-date", value: "12-07-2022"},
	//}, http.StatusOK},
	//{"make-reservation-json", "/reservation-json", "POST", []postData{
	//	{key: "start", value: "13-07-2022"},
	//	{key: "end-", value: "13-07-2022"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	// how are we going to create a web server that is actually going to return status code, something we can post to
	// we need to create a serc=ver and a client as well
	ts := httptest.NewTLSServer(routes) // http test server
	defer ts.Close()

	// table test
	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url) // ts.Url contains the URL server is listening to and, we  need to append it toh url before sending requests
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			//values := url.Values{}
			//for _, x := range e.params {
			//	values.Set(x.key, x.value)
			//}
			//resp, err := ts.Client().PostForm(ts.URL+e.url, values) // ts.Url contains the URL server is listening to and, we  need to append it toh url before sending requests
			//if err != nil {
			//	t.Log(err)
			//	t.Fatal(err)
			//}
			//if resp.StatusCode != e.expectedStatusCode {
			//	t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			//}
		}
	}
}

func TestRepository_MakeReservation(t *testing.T) {
	reservation := models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomId:    1,
		Room: models.Room{
			RoomName: "General's Quarter",
			ID:       1,
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// new recorder httptest.NewRecorder
	rr := httptest.NewRecorder() // it simulates a request response lifecycle - when some hits a browser, hits website , request something and then sends response
	sessionManager.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	// this is important because without this there will be no session to get data,
	// we wont put anything in it to test the case in which session doesn't have reservation
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	sessionManager.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// make table test
func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomId:    1,
		Room: models.Room{
			RoomName: "General's Quarter",
			ID:       1,
		},
	}

	reqBody := "start_date=09-01-2023"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=10-01-2023")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jfon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=xxrr@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=07766909220")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody)) // need to post a body
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // :-) this is important

	// new recorder httptest.NewRecorder
	rr := httptest.NewRecorder() // it simulates a request response lifecycle - when some hits a browser, hits website , request something and then sends response
	sessionManager.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for nothing in the body of request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // :-) this is important

	rr = httptest.NewRecorder()

	sessionManager.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for form invalidation
	reqBody = "start_date=09-01-2023"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=10-01-2023")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jfon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=xxrr@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=077609220")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	sessionManager.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for insert rooms failed
	reqBody = "start_date=09-01-2023"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=10-01-2023")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jfon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Loe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=xxrr@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=07766909220")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // :-) this is important

	rr = httptest.NewRecorder()

	sessionManager.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for insert restriction failed
	reqBody = "start_date=09-01-2023"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=10-01-2023")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jfon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=xxrr@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=07766909220")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // :-) this is important

	rr = httptest.NewRecorder()
	reservation.RoomId = 2
	sessionManager.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returns wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostSearchAvailability(t *testing.T) {

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// easy way to post form values using URL.Values
	//reqBody := "start=09-01-2023"
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "end=10-01-2023")
	//reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData := url.Values{}
	postedData.Add("start", "09-01-2023")
	postedData.Add("end", "10-01-2023")
	postedData.Add("room_id", "1")
	req, _ := http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode())) // need to post a body
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // :-) this is important

	// new recorder httptest.NewRecorder
	rr := httptest.NewRecorder() // it simulates a request response lifecycle - when some hits a browser, hits website , request something and then sends response

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Errorf("PostReservation handler returns wrong response code: got %d", err)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := sessionManager.Load(req.Context(), req.Header.Get("X-Session")) // we need this header to be able to read from session
	if err != nil {
		log.Print(err)
	}
	return ctx
}

// we cant run application using go run cmd/web/*.go command anymore - will have to create script file dor this
//go run: cannot run *_test.go files (cmd/web/main_test.go)
//go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go - use this instead
