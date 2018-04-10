package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//AppUser contains the user information usable by the crud login controller and the db facebook tables
type AppUser struct {
	Valid      bool
	ID         string
	FacebookID string
	FullName   string
	Email      string
}

//structure mimics that of JSON returned by Facebook's API
type graphUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//structure mimics that of JSON returned by Facebook's API
type fbErrorCatcher struct {
	Error fbError `json:"error"`
}

//structure mimics that of JSON returned by Facebook's API
type fbError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FBTraceID string `json:"fbtrace_id"`
}

const graphRequestAddress = "https://graph.facebook.com/v2.3/me"
const fieldsRequested = "fields=id,name,email"

//return user information from facebook. if token was invalid, log error, and return AppUser{Valid:false}
func getFBUser(accesstoken string) AppUser {

	graphRequest := getFacebookAddress(accesstoken)
	//fmt.Println(fbApiUrl) dont print non errors

	response, error := http.Get(graphRequest)
	if error != nil {
		fmt.Println(error)
		return AppUser{Valid: false}
	}
	defer response.Body.Close()

	return *parseAPIResponse(response)
}

//getFacebookAddress returns formatted string to query FB API
func getFacebookAddress(accessToken string) string {
	return graphRequestAddress + "?" + fieldsRequested + "&access_token=" + accessToken
}

//parses API response into
func parseAPIResponse(response *http.Response) *AppUser {

	//read full body
	body, _ := ioutil.ReadAll(response.Body)

	wasError, _ := wasFBError(body)
	if wasError {

		return &AppUser{Valid: false}
	}

	APIResponseObject := parseJSONResponse(body)
	appUser := APIResponseObject.toAppUser()

	return appUser

}

//parses FB API response into a FBAPIObject
func parseJSONResponse(data []byte) (fbObject *graphUser) {
	fbObject = &graphUser{}
	json.Unmarshal(data, &fbObject)

	return
}

//FBAPIObject function to return its appuser version
func (fbObject *graphUser) toAppUser() (businessUser *AppUser) {
	businessUser = &AppUser{Valid: true}

	businessUser.FullName = fbObject.Name
	businessUser.FacebookID = fbObject.ID
	businessUser.Email = fbObject.Email

	return
}

//so every request does not need to make a new empty error for comparison
var emptyError = fbError{}

func wasFBError(data []byte) (bool, fbError) {
	errorCatcher := fbErrorCatcher{}

	json.Unmarshal(data, &errorCatcher)

	return errorCatcher.Error != emptyError, errorCatcher.Error
}

const errorLogFormat = `** Facebook authentication error: 
** message was: %s
** type was: %s
** code was: %d
** facebook trace id was: %s
`

func (error fbError) String() string {
	return fmt.Sprintf(errorLogFormat, error.Message, error.Type, error.Code, error.FBTraceID)
}
