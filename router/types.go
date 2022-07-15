package main

type Data struct {
	EmitterId string `json:emitterId`
	Value     string `json:value`
}

type Command struct {
	TargetId string `json:targetId`
	Value    string `json:value`
}
