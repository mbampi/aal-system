package aal

type Query string

type URI string

type Finding struct {
	Name    string
	Patient string
	Sensor  string
	Value   string
}
