package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type AppUser struct{
	Valid bool
	FirstName string
	FullName string
	Gender string
}

//structure mimics that of JSON returned by Facebook's API
type fbApiObject struct{
	FirstName string `json:"first_name"`
	Name string `json:"name"`
	Gender string `json:"gender"`
}

//structure mimics that of JSON returned by Facebook's API
type fbErrorCatcher struct{
	Error fbError `json:"error"`
}

//structure mimics that of JSON returned by Facebook's API
type fbError struct{
	Message string `json:"message"`
	Type string `json:"type"`
	Code int `json:"code"`
	FBTraceId string `json:"fbtrace_id"`
}



const fbApiRequestBase = "https://graph.facebook.com/v2.3/me"
const fbApifields = "fields=first_name,name,gender,age_range"

//return user information from facebook. if token was invalid, log error, and return AppUser{Valid:false}
func getFBUser(accesstoken string) AppUser {

	fbApiUrl := makeFBApiGetUrl(accesstoken)
	fmt.Println(fbApiUrl)

	response, error := http.Get(fbApiUrl)
	if error != nil{
		fmt.Println(error)
		return AppUser{Valid:false}
	}
	defer response.Body.Close()

	return *parsefbApiResponse(response)
}

func makeFBApiGetUrl(accessToken string) string{
	return fbApiRequestBase + "?" +fbApifields + "&access_token=" +accessToken
}

func parsefbApiResponse(response *http.Response) *AppUser{

	//read full body
	body, _ := ioutil.ReadAll(response.Body)

	wasError, fbError := wasFBError(body)
	if  wasError {
		//log error and return appUser s.t. Valid = false
		fmt.Println(fbError)
		return &AppUser{Valid:false}
	}else{
		fbApiObject := parseJsonToFbApiObject(body)
		appUser := fbApiObject.toAppUser()
		return appUser
	}
}

func parseJsonToFbApiObject(data []byte) (fbObject *fbApiObject){
	fbObject = &fbApiObject{}
	json.Unmarshal(data,&fbObject)
	return
}

func (fbObject *fbApiObject) toAppUser() (businessUser *AppUser){
	businessUser = &AppUser{Valid:true}

	businessUser.FirstName = fbObject.FirstName
	businessUser.FullName = fbObject.Name
	businessUser.Gender = fbObject.Gender

	return
}

//so every request does not need to make a new empty error for comparison
var emptyError = fbError{}

func wasFBError(data []byte) (bool, fbError){
	errorCatcher := fbErrorCatcher{}

	json.Unmarshal(data,&errorCatcher)

	return errorCatcher.Error != emptyError, errorCatcher.Error
}

const errorLogFormat = `** Facebook authentication error: 
** message was: %s
** type was: %s
** code was: %d
** facebook trace id was: %s
`

func (error fbError) String() string{
	return fmt.Sprintf(errorLogFormat,error.Message,error.Type,error.Code,error.FBTraceId)
}

