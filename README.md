# Installing Mongo and getting set up
## Prerequisites:
- Docker Desktop

## Steps
- See [this guide](https://www.thepolyglotdeveloper.com/2019/01/getting-started-mongodb-docker-container-deployment/) for more information.
- `docker pull mongo`
- `docker run -d -p 27017-27019:27017-27019 --name mongodb mongo:latest`
- `docker exec -it mongodb bash`
- `mongo`
- `show dbs`
- `use hungry`
- `db.experiences.save({ name: "Go Vegan" })`
- `db.chefs.save({ name: "Tammy Pensacola" })`
- `db.orders.save({ total: 109.44 })`

# Starting the server
- Place this directory in your $GOPATH: `~/go/src/HUNGRY` on Mac
- `go get github.com/gorilla/mux`
- `go get go.mongodb.org/mongo-driver`
- `go run main.go`

# Endpoints

- All endpoints can be tested in Postman. If there are arguments to send in a POST, you may send them in a raw JSON format body. At least one argument is required in the body of each POST.

## Experiences
- Create: POST to `http://localhost:12345/experience` with:
    - "name" (string)
- Get all: GET to `http://localhost:12345/experiences`
- Get one: GET to `http://localhost:12345/experience/{id}` where `{id}` is an experience ID from the database.

## Chefs
- Create: POST to `http://localhost:12345/chef` with:
    - "date" (time.Time)
    - "ordernumber" (integer) 
    - "chefid" (integer)
    - "experience" (Experience) 
    - "headcount" (integer) 
    - "subtotal" (floating point number)
    - "tax" (floating point number)
    - "tip" (floating point number)
    - "total" (floating point number)
- Get all: GET to `http://localhost:12345/chefs`
- Get one: GET to `http://localhost:12345/chef/{id}` where `{id}` is a chef ID from the database.

## Orders
- Create: POST to `http://localhost:12345/order` with
    - "name" (string)
    - "email" (string)
    - "virtualexperiencesoffered" (Experience array)
- Get all: GET to `http://localhost:12345/orders`
- Get one: GET to `http://localhost:12345/orders/{id}` where `{id}` is an order ID from the database.

Note: I followed the tutorial [here](https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/) for this project.