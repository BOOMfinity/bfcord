package gateway

import "github.com/BOOMfinity/bfcord/gateway/intents"

type Identify struct {
	Properties     IdentifyProperties `json:"properties"`
	Token          string             `json:"token"`
	Shard          []uint16           `json:"shard"`
	Intents        intents.Intent     `json:"intents"`
	LargeThreshold uint16             `json:"large_threshold,omitempty"`
	Compress       bool               `json:"compress,omitempty"`
}

type IdentifyProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}
