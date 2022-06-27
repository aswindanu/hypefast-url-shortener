# URL Shortener
A RESTful API simple ecommerce
with Go using **gorilla/mux** (A nice mux library) and **gorm** (An ORM for Go)

```
API Endpoint : http://localhost:3000
```

## Structure
```
├── .github
│   └── workflows
│       └── main.yml	 		 	 // Github Actions CI (Go test runner)
├── app
│   ├── app.go
│   ├── handler          		 	 // Our API core handlers
│   │   ├── urlshortener  		 	 // APIs for urlshortener model
│   │   └── common.go    		 	 // Common response functions
│   └── model
│       └── model.go	 		 	 // Models for application
├── config
│   └── config.go 	 		 	 // Configuration
└── main.go
└── Dockerfile
└── docker-compose.yml
└── hypefast.postman_collection.json 	 	 // Postman collection API
```

## How To Run

1. Clone this repo

	```
	git clone https://github.com/aswindanu/hypefast-url-shortener.git
	```

2. Go into directory, the build and run docker-compose using this command

	```
	cd hypefast-url-shortener
	docker-compose up -d --build
	```

3. API will exposed on this url `http://localhost:3000`, 

#### Note
- Database using PostgreSQL, can be access using this credentials 

	```
	HOST = 127.0.0.1
	PORT = 5454
	DATABASE NAME = todoapp
	USERNAME = root
	PASSWORD = password
	```

- Please test the API for details using Postman collection that can imported from this file `hypefast.postman_collection.json`

- Result of `github actions` can be seen from tab `Action` on github repo itself

## API HIGHLIGHT 
###### *) please import hypefast.postman_collection.json for more details

#### /api/v1/url-shortener
* `GET` : Get all urls
* `POST` : Add new url
	
		body: 
		{
			"original_url": "<url>"
		}

#### /api/v1/url-shortener/{uid}
* `GET` : Get detail url

#### /api/v1/url-shortener/redirect
* `GET` : Redirect to original_url

		Query params:
		- alternate_url : <alternate_url>

#### /api/v1/url-shortener/{uid}/activate
* `POST` : Activate short url

#### /api/v1/url-shortener/{uid}/deactivate
* `POST` : Deactivate short url


### EXPECTED SCENARIO

1. Generate new short url with POST `/api/v1/url-shortener` endpoint
2. Check details short url with GET `/api/v1/url-shortener/{uid}` endpoint
3. Redirect short url with GET `/api/v1/url-shortener/redirect` endpoint and query params `alternate_url`
4. Deactivate url with POST `/api/v1/url-shortener/{uid}/deactivate` endpoint
5. Try step 2 & 3 again
6. Re-activate url with POST `/api/v1/url-shortener/{uid}/activate` endpoint
7. Try step 2 & 3 again
