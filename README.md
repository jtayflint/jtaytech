# JTAYTECH Github Project
This project proivdes code to share to the general public.  I have add a MIT license to this code to grant others the right to copy and use this code.

Currently, this project only has 1 application. The application is weatherapi.

## weatherapi
weatherapi is a sample app created for a techcnical interview. See the App Requests below for more informaiton.

### Project Submission - Weather Service Assignment 
Write an HTTP server that serves the current weather. Your server should expose an endpoint that: 
1. Accepts latitude and longitude coordinates 
2. Returns the short forecast for that area for Today (“Partly Cloudy” etc) 
3. Returns a characterization of whether the temperature is “hot”, “cold”, or “moderate” (use your discretion on mapping temperatures to each type) 
4. Use the [National Weather Service API Web Service](https://urldefense.proofpoint.com/v2/url?u=https-3A__www.weather.gov_documentation_services-2Dweb-2Dapi&d=DwMFaQ&c=_EdSgJoS8igo01XnekBu_azVXoUPxJkwz9O2AzwhBbE&r=rhZn4o-SVBKWg8RToGhlXOJm2K5NDseuaAY8rdkTwIs&m=z42ii4Dc2Ny7rIqIyR8EpNpn6t5jWKW5nQkKeMAUvujP7Ul4RujzeKLAd4PbL97D&s=QX-wbegcMSCQI15AUvxgT7_HaM5fy7RzNtgFYdtsRoQ&e=) as a data source. 
 
The purpose of this exercise is to provide a sample of your work that we can discuss together in the Technical Interview. 
We respect your time. Spend as long as you need, but we intend it to take around an hour. 
We do not expect a production-ready service, but you might want to comment on your shortcuts. 
The submitted project should build and have brief instructions so we can verify that it works. 
You may write in whatever language or stack you're most comfortable in, but it's recommended to use the language for the job you're applying for (Go).

### Initialize Project Commands
```bash
mkdir weatherapi
cd weatherapi
go mod init weatherapi
go get github.com/go-chi/chi/v5
go get github.com/caarlos0/env/v11
go get go.uber.org/fx
go get go.uber.org/zap
```

### Run Application Commands


```bash
go mod tidy
go build .
go  run .
```

### Sample Request
http://localhost:8080/weather/lat/38.8062/lon/-94.8368
