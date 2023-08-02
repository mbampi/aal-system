package aal

import "time"

type Query string

type URI string

type Finding struct {
	Name      string    `json:"name"`
	Patient   string    `json:"patient"`
	Sensor    string    `json:"sensor"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
