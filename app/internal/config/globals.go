// Package config constants (mostly keys) used for accessing data in a centralized manner
//
// This package should be as independent as possible to avoid circular dependencies
package config

// route paths for gin
const (
	RouteLogin    = "/login"
	RouteLogout   = "/logout"
	RouteRegister = "/register"
)
