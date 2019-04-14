# Server for the Valuator
This application creates a rudimentary web server that uses the valuator
package and provides a simple query interface for the user to get filing data,
valuation metrics and the valuation of a ticker based on trends.

The purpose of the application is to demonstrate the use of the valuator
package. It uses HTML templates to format the data from the valuator to
create a page that serves the information provided by the package to the user.

Usage:

- Build the package
ex: go build
- ./server

The server runs by default on localhost:8080 and should be accesible via a web
browser.
