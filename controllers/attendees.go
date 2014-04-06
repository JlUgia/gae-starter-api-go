package controllers

import (
        "net/http"
        "io"
        "time"
        "errors"
        "encoding/json"

        "appengine"
        "appengine/datastore"
)

type Attendee struct {
    Name  		string     `json:"name"`
    Email 		string     `json:"email"`
    CreatedAt   time.Time  `json:"created_at"`
}

func init() {
        http.HandleFunc("/users", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {

    c := appengine.NewContext(r)

    val, err := route(c, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response, error := json.MarshalIndent(val, "", "    ")
    if(error != nil) {
        http.Error(w, error.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(response)
}

func route(c appengine.Context, r *http.Request) (interface{}, error) {

    switch r.Method {

    case "POST":
        attendee, err := decodeObject(r.Body)
        if err != nil {
            return nil, err
        }
        return registerUser(c, attendee)

    case "GET":
        return getAllUsers(c)
    }

    return nil, errors.New("Method not implemented")
}

func decodeObject(body io.ReadCloser) (*Attendee, error) {

    defer body.Close()

    var attendee Attendee
    err := json.NewDecoder(body).Decode(&attendee)
    return &attendee, err
}

func registerUser(c appengine.Context, attendee *Attendee) (*Attendee, error) {

        attendee.CreatedAt = time.Now()
        _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "attendee", nil), attendee)
        if err != nil {
            return nil, err
        }

        return attendee, nil
}

func getAllUsers(c appengine.Context) ([]Attendee, error) {

        attendes := []Attendee{}

        _, err := datastore.NewQuery("attendee").Order("CreatedAt").GetAll(c, &attendes)
        if err != nil {
            return nil, err
        }

        return attendes, nil
}