package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Experience struct {
	ID	 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string				`json:"name,omitempty" bson:"name,omitempty"`
}

type Order struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date           time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Ordernumber    int                `json:"ordernumber,omitempty" bson:"ordernumber,omitempty"`
	Chefid         int				  `json:"chefid,omitempty" bson:"chefid,omitempty"`
	Experience     *Experience    	  `json:"experience,omitempty" bson:"experience,omitempty"`
	Headcount      int				  `json:"headcount,omitempty" bson:"headcount,omitempty"`
	Subtotal       float32			  `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
	Tax            float32			  `json:"tax,omitempty" bson:"tax,omitempty"`
	Tip            float32			  `json:"tip,omitempty" bson:"tip,omitempty"`
	Total          float32			  `json:"total,omitempty" bson:"total,omitempty"`
}

type Chef struct {
	ID				          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name			          string			 `json:"name,omitempty" bson:"name,omitempty"`
	Email			          string			 `json:"email,omitempty" bson:"email,omitempty"`
	Virtualexperiencesoffered []Experience		 `json:"virtualexperiencesoffered,omitempty" bson:"virtualexperiencesoffered,omitempty"`
}

// Experience endpoints
func CreateExperienceEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var experience Experience
	_ = json.NewDecoder(request.Body).Decode(&experience)
	collection := client.Database("hungry").Collection("experiences")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, experience)
	json.NewEncoder(response).Encode(result)
}

func GetExperienceEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var experience Experience
	collection := client.Database("hungry").Collection("experiences")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Experience{ID: id}).Decode(&experience)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(experience)
}

func GetExperiencesEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var experiences []Experience
	collection := client.Database("hungry").Collection("experiences")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var experience Experience
		cursor.Decode(&experience)
		experiences = append(experiences, experience)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(experiences)
}

// Order endpoints
func CreateOrderEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var order Order
	_ = json.NewDecoder(request.Body).Decode(&order)

	// validation
	if order.Total < 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "Order total must be greater than or equal to $0.00" }`))
		return
	}

	if order.Date.IsZero() {
		order.Date = time.Now()
	}

	collection := client.Database("hungry").Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, order)
	json.NewEncoder(response).Encode(result)
}

func GetOrderEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var order Order
	collection := client.Database("hungry").Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Order{ID: id}).Decode(&order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(order)
}

func GetOrdersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var orders []Order
	collection := client.Database("hungry").Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var order Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(orders)
}

// Chef endpoints
func CreateChefEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var chef Chef
	_ = json.NewDecoder(request.Body).Decode(&chef)
	collection := client.Database("hungry").Collection("chefs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, chef)
	json.NewEncoder(response).Encode(result)
}

func GetChefEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var chef Chef
	collection := client.Database("hungry").Collection("chefs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Chef{ID: id}).Decode(&chef)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(chef)
}

func GetChefsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var chefs []Chef
	collection := client.Database("hungry").Collection("chefs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var chef Chef
		cursor.Decode(&chef)
		chefs = append(chefs, chef)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(chefs)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	router.HandleFunc("/experience", CreateExperienceEndpoint).Methods("POST")
	router.HandleFunc("/experience/{id}", GetExperienceEndpoint).Methods("GET")
	router.HandleFunc("/experiences", GetExperiencesEndpoint).Methods("GET")
	
	router.HandleFunc("/order", CreateOrderEndpoint).Methods("POST")
	router.HandleFunc("/order/{id}", GetOrderEndpoint).Methods("GET")
	router.HandleFunc("/orders", GetOrdersEndpoint).Methods("GET")
	
	router.HandleFunc("/chef", CreateChefEndpoint).Methods("POST")
	router.HandleFunc("/chef/{id}", GetChefEndpoint).Methods("GET")
	router.HandleFunc("/chefs", GetChefsEndpoint).Methods("GET")
	
	http.ListenAndServe(":12345", router)
}
