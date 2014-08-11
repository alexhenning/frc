// dsmutex acts as a mutex for the FRC driverstation for testing
// purposes.
//
// Usage:
//     teststand serve
//     teststand client [--name=me]
//     teststand take [--name=me]              Take the mutex. Prints status while waiting.
//     teststand give [--name=me]              Give the mutex to someone waiting.
//     teststand message [--name=me] "Message" Message others waiting on the mutex know.
//     teststand ds [--name=me] enable         Set the ds state to either enable or disable.
//
// Using wget as an API:
//
//     Taking: wget -q --read-timeout=0 10.1.90.2:8080/take?name=$USER -O -
//     Giving: wget -q 10.1.90.2:8080/give?name=$USER -O -
//     Msging: wget -q "10.1.90.2:8080/give?name=$USER&Message=Hi all!" -O -
package main
