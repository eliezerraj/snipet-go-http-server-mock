# snipet-go-http-server-mock

Server http to mock behavior

How to use

1. Create a go module
go mod init github/snipet-go-http-server-mock/main

2. Create a .env with name of POD and http PORT
PORT=:8900
NAME_POD=backend-01

3. Run the server
go run .

4. Check all internal parameters
GET http://localhost:8900/

5. Check all internal parameters
POST http://localhost:8900/setup
	{
		"response_time":4,
		"response_status_code":200,
		"is_random_time": true
	}
response_time = set the time to response any request in SECONDS
response_status_code = set the status code to any request
is_random_time (true) set a random time response since 0 to response_time

6. Return 1 payload with a customer fake data
GET http://localhost:8900/customer_fake

6. Return a payload with a list with 50 customer fake
GET http://localhost:8900/list_customer_fake

6. Post a number to calc a Fibonacci and stress the CPU
POST http://localhost:8900/stress_cpu
		{
			"count":200
		}
count = number to calc and stress CPU

7. Create a Docker image
docker build -t go_http_server_mock .

8. Run a docker file
docker run -dit --name go_http_server_mock -p 8900:8900 go_http_server_mock

9. Create a Docker Builder image (light image)
docker build -t go_http_server_mock . -f Dockerfile-Builder

9. Run a docker Builder image
docker run -dit --name go_http_server_mock -p 8900:8900 go_http_server_mock 