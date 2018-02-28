package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type fbUser struct {
	valid     bool
	firstName string
	lastNAme  string
	//any other information available through facebooks graph api
}

const fbApiRequestBase = "https://graph.facebook.com/v2.3/me"
const fbApifields = "fields=first_name,gender,age_range"

func getFBUser(accesstoken string) fbUser {

	fbApiUrl := makeFBApiGetUrl(accesstoken)
	fmt.Println(fbApiUrl)

	response, error := http.Get(fbApiUrl)
	if error != nil {
		fmt.Println(error)
		return fbUser{valid: false}
	}
	defer response.Body.Close()

	return parsefbApiResponse(response)
}

func makeFBApiGetUrl(accessToken string) string {
	return fbApiRequestBase + "?" + fbApifields + "&access_token=" + accessToken
}

func parsefbApiResponse(response *http.Response) fbUser {
	user := fbUser{}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	//todo: parse json into business object

	return user
}
