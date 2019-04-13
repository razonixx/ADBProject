package main

import "gopkg.in/mgo.v2"

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	sentiment "gopkg.in/vmarkovtsev/BiDiSentiment.v1"
)

type Review struct {
	Date         string
	DateAdded    string
	DateSeen     string
	DidPurchase  string
	DoRecommend  string
	ID           string
	NumHelpful   int
	Rating       int
	SourceURLs   string
	Text         string
	Title        string
	UserCity     string
	UserProvince string
	Username     string
}

type ReviewData struct {
	ID           string
	Name         string
	Asins        string
	Brand        string
	Categories   string
	Keys         string
	Manufacturer string
	Reviews      Review
}

type Result struct {
	ProductName string
	Review      string
	Sentiment   float32
}

func main() {
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Connect to local MongoDB
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	session.SetMode(mgo.Monotonic, true)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	collectionRead := session.DB("ADB").C("Amazon")
	collectionWrite := client.Database("ADB").Collection("AmazonOutput")
	fmt.Println("Connected to MongoDB!")
	//////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Iterate through whole collection and perform sentiment analysis
	//When done, write the results into our new Mongo Collection
	find := collectionRead.Find(bson.M{})

	var review ReviewData

	items := find.Iter()
	i := 0
	arrayOfTexts := []string{}
	arrayOfReviews := []ReviewData{}
	sentimentSession, _ := sentiment.OpenSession()
	for items.Next(&review) {
		arrayOfTexts = append(arrayOfTexts, review.Reviews.Text)
		arrayOfReviews = append(arrayOfReviews, review)
		i++
	}
	sentimentFound, err := sentiment.Evaluate(arrayOfTexts, sentimentSession)
	for j := 0; j < i; j++ {
		tempResult := Result{arrayOfReviews[j].Name, arrayOfTexts[j], sentimentFound[j]}
		fmt.Printf("Name: %s Sentiment %f\n", tempResult.ProductName, tempResult.Sentiment)
		_, err := collectionWrite.InsertOne(context.TODO(), tempResult)
		if err != nil {
			log.Fatal(err)
		}
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Close connection with MongoDB

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	session.Close()
	fmt.Println("Connection to MongoDB closed.")
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
