package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"
)

func main() {
	//token := ""
	token := os.Args[1]
	//fmt.Println(token)

	client, err := elastic.NewClient(
		elastic.SetURL("https://es.int.panda.dpccloud.com"),
		elastic.SetBasicAuth("cpsiam", "f53b0cdcf9e47f030432019382ab9e545008eaaa"),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false))
	if err != nil {
		log.Println("could not connect...", err)
		return
	}

	main1(client, token)
	main2(client, token)

}

func main1(client *elastic.Client, token string) {
	defer timeTrack(time.Now(), "main1()")
	//iam:tokens/access_token/177e5be3-485b-48e3-9129-860f37a10443
	//token := "60b0f54f-9d9b-4a53-95b4-a9b49e38582e"
	callES("iam:tokens", "access_token", token, "", client)
	agentID := "d4a2784b-92e5-49ee-8a7b-347ded7e257b"
	agentTenantID := "85b281cf-c074-4d08-80d5-4fd98458640f"
	callES("iam:identities", "agent", agentID, agentTenantID, client)
	tenantID := "23223096-68e6-44ad-b1b2-183fbf489090"
	callES("iam:identities", "tenant", tenantID, "", client)

	userID := "4422b6ed-7f26-491a-b04f-e9ddb8094221"
	callES("iam:identities", "user", userID, tenantID, client)

}

func main2(client *elastic.Client, token string) {
	defer timeTrack(time.Now(), "main2()")
	agentID := "d4a2784b-92e5-49ee-8a7b-347ded7e257b"
	agentTenantID := "85b281cf-c074-4d08-80d5-4fd98458640f"
	tenantID := "23223096-68e6-44ad-b1b2-183fbf489090"
	userID := "4422b6ed-7f26-491a-b04f-e9ddb8094221"

	callES("iam:tokens", "access_token", token, "", client)

	res, err := client.MultiGet().
		Add(elastic.NewMultiGetItem().Index("iam:identities").Type("agent").Id(agentID).Routing(agentTenantID)).
		Add(elastic.NewMultiGetItem().Index("iam:identities").Type("tenant").Id(tenantID)).
		Add(elastic.NewMultiGetItem().Index("iam:identities").Type("user").Id(userID).Routing(tenantID)).Do()
	if err != nil {
		fmt.Println(err)
	}

	/* Decode the json into map of string to interfaces */
	for index := 0; index < len(res.Docs); index++ {
		doc := make(map[string]interface{})
		r := bytes.NewReader(*res.Docs[index].Source)
		decoder := json.NewDecoder(r)
		err = decoder.Decode(&doc)

		fmt.Println(doc)

	}

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n%s took %s", name, elapsed)
}

/**
 * Function to get the document by id
 *
 * indexName                  : Name of index
 * typeName                   : Type of document
 * entityId                   : Document id
 * routingID                  : Routing ID for parent-child relationship
 *
 * returns map of string to interfaces, error if fails
 */
func callES(indexName, typeName, entityID, routingID string, client *elastic.Client) (
	searchResult map[string]interface{}, err error) {

	/* Check for proper arguments */
	if len(indexName) == 0 || len(typeName) == 0 || len(entityID) == 0 {
		log.Println("Bad args...")
		return
	}

	/* Get the document with given id */
	get, err := client.Get().
		Index(indexName).
		Type(typeName).
		Id(entityID).
		Routing(routingID).
		Do()

	// If error, log it and return
	if err != nil {
		log.Println("Error searching...", err)
		return
	}

	/* Check whether we found the document */
	if !get.Found {
		log.Println("Error Found...", get.Found)
		return
	}

	log.Println("Found", get.Found)

	/* Decode the json into map of string to interfaces */
	r := bytes.NewReader(*get.Source)
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&searchResult)

	// If error, log it and return
	if err != nil {
		log.Println("Error Marshalling...", err)
		return
	}

	return
}
