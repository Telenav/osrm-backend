// Package searchconnector provides api to communicate with search server.
// It starts during initialization and stops when calling Stop().
//
// When user make a call to ChargeStationSearch
// - User will get a channel of ChargeStationsResponse immediately
// - User continues his logic, he will be blocked if trying to fetch response
//   from channel before the response is ready.
// - TNSearchConnector will communicate with TN search server via http and assemble the
//   response, when everything is finished he will put result into the response
//   channel.  If request failed or timeout, he will put related error in the
//   response.
// - User could get response and error information from channel
//
// This package wraps the logic of slow IO communication, provides flexibility
// to external user for when to wait for the response.
package searchconnector
