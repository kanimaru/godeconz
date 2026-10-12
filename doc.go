// Package godeconz holds the data model of the deCONZ REST API (lights, sensors, groups, scenes,
// gateway configuration) shared by the REST client in package rest and the websocket client in
// package ws.
//
// Read models use pointer fields for values a device may not report, so "not reported" and
// "zero" stay distinguishable. The same models decode websocket "changed" events, which carry
// only the attributes that changed. Update types (LightStateUpdate, GroupActionUpdate, ...)
// send only the fields that are set.
//
// deCONZ documentation: https://dresden-elektronik.github.io/deconz-rest-doc/
package godeconz

// Ptr returns a pointer to v. It keeps update literals short:
//
//	godeconz.LightStateUpdate{On: godeconz.Ptr(true), Bri: godeconz.Ptr[uint8](128)}
func Ptr[T any](v T) *T {
	return &v
}
