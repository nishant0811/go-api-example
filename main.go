package main

import (
    "context"
    "crypto/md5"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)


var collection , postCollection = ConnecttoDB()



func main() {
  http.HandleFunc("/user" , createUser)
  http.HandleFunc("/user/" , findUser)
  http.HandleFunc("/posts" , createPost)
  http.HandleFunc("/posts/" , findPost)
  http.HandleFunc("/posts/users/" , userPost)
	// set our port address as 8081
	log.Fatal(http.ListenAndServe(":8081", nil))
}


func ConnecttoDB() (*mongo.Collection , *mongo.Collection ){

	// Set client options
	//change the URI according to your database
	clientOptions := options.Client().ApplyURI("mongodb+srv://Nishant:nishant1234@cluster0.m0yjk.mongodb.net/hospital?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//DB collection address which we are going to use
	//available to functions of all scope
	collection := client.Database("Appointy").Collection("NewsUsers")
  postCollection := client.Database("Appointy").Collection("NewPosts")
	return collection, postCollection
}

//Function to create a new User in Database
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

  if r.Method == "POST"{
    r.ParseForm()

    pass := md5.Sum([]byte(r.Form.Get("email")));
	  result, err := collection.InsertOne(context.TODO(),bson.D{
    {Key: "name", Value: r.Form.Get("name")},
    {Key: "passowrd", Value: string(pass[:])},
    {Key : "email" , Value : r.Form.Get("email")},
    {Key : "userId" , Value : r.Form.Get("id")},
})

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result)
  } else {
    fmt.Fprintf(w, "Only Post Request Accepted")
			return
  }

}


func findUser(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  keys,ok := r.URL.Query()["id"];
  if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'id' is missing")
        return
    }
  key := keys[0];

  var result bson.M
  fmt.Println(key)
  err := collection.FindOne(context.TODO() , bson.M{"userId" : key}).Decode(&result)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(result)
  json.NewEncoder(w).Encode(result)
}

func createPost(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")

    if r.Method == "POST"{
      r.ParseForm()

      dt := time.Now();
      dt.Format("01-02-2006 15:04:05")
  	  result, err := postCollection.InsertOne(context.TODO(),bson.D{
      {Key: "id", Value: r.Form.Get("pid")},
      {Key: "caption", Value: r.Form.Get("caption")},
      {Key : "url" , Value : r.Form.Get("url")},
      {Key : "time" , Value : dt.String()},
      {Key : "userId" , Value : r.Form.Get("userId")},
  })

  	if err != nil {
  		log.Fatal(err)
  	}

  	json.NewEncoder(w).Encode(result)
    } else {
      fmt.Fprintf(w, "Only Post Request Accepted")
  			return
    }

}


func findPost(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  keys,ok := r.URL.Query()["id"];
  if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'id' is missing")
        return
    }
  key := keys[0];

  var result bson.M
  fmt.Println(key)
  err := postCollection.FindOne(context.TODO() , bson.M{"id" : key}).Decode(&result)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(result)
  json.NewEncoder(w).Encode(result)
}


func userPost(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  keys,ok := r.URL.Query()["id"];
  if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'id' is missing")
        return
    }
  key := keys[0];

  filterCursor, err := postCollection.Find(context.TODO(), bson.M{"userId": key})
  if err != nil {
      log.Fatal(err)
  }

  var postFilter []bson.M
  if err = filterCursor.All(context.TODO(), &postFilter); err != nil {
      log.Fatal(err)
  }

  json.NewEncoder(w).Encode(postFilter)

}
