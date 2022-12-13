# Analog API Backend

This is the webserver for the Analog Restful API written with Golang. The objective is to provide a robust & performant server connected to a PostgreSQL database to serve clients.

## Getting started

You must have at least :

* Golang installed
* PostgreSQL running

1. Defining required environment variable

First of all, please copy file `.env-example`, remove each comment and replace value with your own secrets.

2. Run webserver (default port `8080`)

Execute command `go run .`, you should have on STDOUT a URL with the running application.

3. Display all cameras

`curl http://localhost:8008/camera/`, you should have a list of cameras (you must have added cameras in database)