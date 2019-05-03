Sentiment Analysis applied to Dataset using Tensorflow
=========================

- The objective of this project is to analyze a dataset to obtain conclusions using said dataset.
- We decided to analyze a set of reviews left on products on [Amazon](amazon.com). 
- The dataset was obtained from [Kaggle](kaggle.com), and can be found [here](https://www.kaggle.com/datafiniti/consumer-reviews-of-amazon-products#1429_1.csv)
- In order to run the code, the user must have access to a mongodb server, which contains the Kaggle dataset. 
- The mongodb driver for Go must be installed. Instructions on how to do so can be found [here](https://github.com/mongodb/mongo-go-driver). The BiDiSentiment library must also be installed, see [here](https://github.com/vmarkovtsev/BiDiSentiment) for instructions on how to do so.D
- BiDiSentiment uses [Tensorflow](https://www.tensorflow.org/) to implement its sentiment analysis algorithm. As such, Tensorflow must also be installed. Instructions for installing the Tensorflow Go API are found [here](https://www.tensorflow.org/install/lang_go)

- To import the Kaggle dataset into MongoDB, download the .csv found [here](https://www.kaggle.com/datafiniti/consumer-reviews-of-amazon-products/downloads/consumer-reviews-of-amazon-products.zip/4) and use the [mongoimport](https://docs.mongodb.com/manual/reference/program/mongoimport/) tool to create the collection inside a MongoDB database called `ADB`. The dataset must be in a collection called `Amazon`. The output of the program is a new collection inside the same DB called `AmazonOutput`

- Once the dataset is in the local MongoDB database, the user just needs to run `go run sentimentAnalysis.go` if you want to see feedback or `go run sentimentAnalysisOptimized.go` if you want to run an optimized version without feedback.