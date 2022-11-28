package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/quentinchampenois/analog_api"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a main.App

func testMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE FROM cameras")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = a.DB.Exec("ALTER SEQUENCE cameras_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatalln(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS cameras
(
    id SERIAL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    focus TEXT NOT NULL,
    film INTEGER NOT NULL,
    CONSTRAINT cameras_pkey PRIMARY KEY (id)
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/cameras", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentCamera(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/camera/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Camera not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Camera not found'. Got '%s'", m["error"])
	}
}

func TestCreateCamera(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"Minolta AF-S", "type": "Compact Point & Shoot", "focus": "Autofocus", "film": 0}`)
	req, _ := http.NewRequest("POST", "/camera", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Minolta AF-S" {
		t.Errorf("Expected camera name to be 'Minolta AF-S'. Got '%v'", m["name"])
	}

	if m["type"] != "Compact Point & Shoot" {
		t.Errorf("Expected camera type to be 'Compact Point & Shoot'. Got '%v'", m["type"])
	}

	if m["focus"] != "Autofocus" {
		t.Errorf("Expected camera focus to be 'Autofocus'. Got '%v'", m["focus"])
	}

	if m["film"] != 0 {
		t.Errorf("Expected camera film to be '0'. Got '%v'", m["film"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected camera ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetCamera(t *testing.T) {
	clearTable()
	addCameras(1)

	req, _ := http.NewRequest("GET", "/camera/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCameras(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO cameras(name, type, focus, film) VALUES($1, $2, $3, $4)", "Camera "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func testUpdateCamera(t *testing.T) {
	clearTable()
	addCameras(1)

	req, _ := http.NewRequest("GET", "/camera/1", nil)
	response := executeRequest(req)

	var camera map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &camera)
	var updatedBytes = []byte(`{"name":"Fujica DL-100", "type": "Compact Point & Shoot", "focus": "Autofocus", "film": 0}`)
	req, _ = http.NewRequest("PUT", "/camera/1", bytes.NewBuffer(updatedBytes))
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != camera["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", camera["id"], m["id"])
	}

	if m["name"] == camera["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", camera["name"], m["name"], m["name"])
	}

	if m["type"] != camera["type"] {
		t.Errorf("Expected the type to remain the same (%v). Got %v", camera["type"], m["type"])
	}

	if m["focus"] != camera["focus"] {
		t.Errorf("Expected the focus to remain the same (%v). Got %v", camera["focus"], m["focus"])
	}

	if m["film"] != camera["film"] {
		t.Errorf("Expected the film to remain the same (%v). Got %v", camera["film"], m["film"])
	}
}

func TestDeleteCamera(t *testing.T) {
	clearTable()
	addCameras(1)

	req, _ := http.NewRequest("GET", "/camera/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/camera/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/camera/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
