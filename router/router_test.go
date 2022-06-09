package router

import (
	"capstoneProyect/ports"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert" // add Testify package
)

// Structure for specifying input and output data
// of a single test case
type tests struct {
	description  string // description of test case
	route        string // route to test
	expectedCode int    // expected HTTP status code
}

var app = fiber.New()

func TestIndexRoute(t *testing.T) {

	// First test
	test := tests{
		description:  "get HTTP status 200",
		route:        "/",
		expectedCode: 200,
	}

	// route for test
	app.Get("/", ports.Index)

	req := httptest.NewRequest("GET", test.route, nil)

	resp, _ := app.Test(req, 10)

	// check if status code is as expected
	assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

}

func TestGetAllRoute(t *testing.T) {
	// First test
	test := tests{
		description:  "get HTTP status 200",
		route:        "/pokemons",
		expectedCode: 200,
	}

	// Create route with GET method for test
	app.Get("/pokemons", ports.GetAll)

	req := httptest.NewRequest("GET", test.route, nil)

	resp, _ := app.Test(req, 10)
	// fmt.Println(resp)
	// check if status code is as expected
	assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

}

func TestGetIdRoute(t *testing.T) {

	tests := []tests{
		// First test case
		{
			description:  "get HTTP status 200",
			route:        "/pokemons/1",
			expectedCode: 200,
		},
		// Second test case
		{
			description:  "get HTTP status 404, when doesn´t exist",
			route:        "/pokemons/200000",
			expectedCode: 404,
		},
		//Third test case
		{
			description:  "get HTTP status 404, when doesn´t exist",
			route:        "/pokemons/randomString",
			expectedCode: 404,
		},
		//Fourth test case
		{
			description:  "get HTTP status 404, when doesn´t exist",
			route:        "/pokemons/1a2b3c",
			expectedCode: 404,
		},
	}

	// Create route with GET method for test
	app.Get("/pokemons/:id", ports.GetById)

	// Iterate through single test cases
	for _, test := range tests {

		req := httptest.NewRequest("GET", test.route, nil)

		resp, _ := app.Test(req, 10)

		// check if status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
