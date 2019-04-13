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
	Rating      int
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
	//model, err := sentiment.Restore()
	if err != nil {
		log.Fatal(err.Error())
	}

	items := find.Iter()
	for items.Next(&review) {
		//analysis := model.SentimentAnalysis(review.Reviews.Text, sentiment.English) // 1 = positive, 0 = negative
		session, _ := sentiment.OpenSession()
		defer session.Close()
		sentimentFound, err := sentiment.Evaluate([]string{review.Reviews.Text}, session)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Review:\n\t%s\nRating:\n\t%d\nSentiment: \n\t%f\n", review.Reviews.Text, review.Reviews.Rating, sentimentFound[0])

		result := Result{review.Name, review.Reviews.Text, review.Reviews.Rating, sentimentFound[0]}
		_, err = collectionWrite.InsertOne(context.TODO(), result)
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
