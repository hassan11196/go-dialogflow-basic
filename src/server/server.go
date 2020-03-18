package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	// dialogflow "cloud.google.com/go/dialogflow/apiv2"
	// dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	// "errors"
)

// only needed below for sample processing
// func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
// 	ctx := context.Background()

// 	sessionClient, err := dialogflow.NewSessionsClient(ctx)
// 	if err != nil {
// 			return "", err
// 	}
// 	defer sessionClient.Close()

// 	if projectID == "" || sessionID == "" {
// 			return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
// 	}

// 	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
// 	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
// 	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
// 	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
// 	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

// 	response, err := sessionClient.DetectIntent(ctx, &request)
// 	if err != nil {
// 			return "", err
// 	}

// 	queryResult := response.GetQueryResult()
// 	fulfillmentText := queryResult.GetFulfillmentText()
// 	return fulfillmentText, nil
// }

func main() {
	os.Setenv("ACCESS_TOKEN", "dd1fdad235274ef7b54ebe454c910a6c")
	os.Setenv("PROJECT_ID", "gdgtest-c33ac")
	fmt.Println("Launching server...")
	baseUrl := "https://api.dialogflow.com/v1/query"

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8085")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		fmt.Print("New messaage:", string(newmessage))
		// send new string back to client
		// conn.Write([]byte(newmessage + "\n"))

		// response,_ := DetectIntentText(os.Getenv("PROJECT_ID"),121542, string(message),"en-US")
		// fmt.Println(response)
		var jsonStr = []byte(`{
			
			"lang": "en",
			"query": ` + string(message) + `,
			"sessionId": "12345",
			"timezone": "America/New_York",
			"v":"20150910"
		  }`)
		req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)

		fmt.Println("response Headers:", resp.Header["Cache-Control"])

		var result map[string]map[string]map[string][]map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Answer:", body)
		res := result["result"]
		ful := res["fulfillment"]
		mess := ful["messages"]
		first_mes := mess[0]
		ful2 := first_mes["speech"]
		fmt.Println(ful2)
		conn.Write([]byte(string(string(ful2) + "\n")))

	}
}
