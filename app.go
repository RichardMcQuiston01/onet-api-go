package main

import (
	"encoding/base64"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func goDotEnvVariable(envFilename string, key string) string {

	// load .env file
	err := godotenv.Load(envFilename)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	var environmentFile = ".env.develop"
	godotenv.Load(environmentFile)

	onetUsername := goDotEnvVariable(environmentFile, "USERNAME")
	onetPassword := goDotEnvVariable(environmentFile, "PASSWORD")
	req.Header.Add("Authorization", "Basic "+basicAuth(onetUsername, onetPassword))
	return nil
}

func main() {
	var environmentFile = ".env.develop"
	godotenv.Load(environmentFile)

	onetUsername := goDotEnvVariable(environmentFile, "USERNAME")
	onetPassword := goDotEnvVariable(environmentFile, "PASSWORD")
	//baseUrl := goDotEnvVariable(environmentFile, "BASEURL")

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", "https://services.onetcenter.org/ws/mnm/interestprofiler/job_zones", nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(onetUsername, onetPassword))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	log.Print(bodyText)
}
