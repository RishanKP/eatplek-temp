package db

import ("context"
         "go.mongodb.org/mongo-driver/mongo"
         "go.mongodb.org/mongo-driver/mongo/options"
         "fmt"
         "log"
         "os"
)


var URI = fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/%s?retryWrites=true&w=majority",os.Getenv("DB_USER"),os.Getenv("DB_PASS"),os.Getenv("DB_CLUSTER"),os.Getenv("DB_NAME"))
var clientOptions = options.Client().ApplyURI(URI)
var Client, Err = mongo.Connect(context.TODO(), clientOptions)

func Connect()  {
		if Err != nil {
		    log.Fatal(Err)
	}

	Err = Client.Ping(context.TODO(), nil)

	if Err != nil {
	    log.Fatal(Err)
	}
}
