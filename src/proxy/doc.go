// Package proxy defines and implements AppGateway: the interface between Alitas
// and an application.
//
// Alitas communicates with the App through an AppGateway interface, which has two
// implementations:
//
// - SocketProxy: A SocketProxy connects to an App via TCP sockets. It enables
// the application to run in a separate process or machine, and to be written in
// any programming language.
//
// - InmemProxy: An InmemProxy uses native callback handlers to integrate Alitas
// as a regular Go dependency.
package proxy
