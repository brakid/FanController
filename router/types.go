package main

type Data struct {
	EmitterId string  `json:emitterId`
	Value     float64 `json:value`
}

type Command struct {
	TargetId string  `json:targetId`
	Value    float64 `json:value`
}

type ComparatorType int

const (
	LOWER ComparatorType = iota
	LOWER_EQUAL
	EQUAL
	GREATER_EQUAL
	GREATER
)

type TransformationFunctionDescriptor struct {
	Comparator  ComparatorType `json:comparator`
	Threshold   float64        `json:threshold`
	TargetIds   []string       `json:targetIds`
	ResultTrue  float64        `json:resultTrue`
	ResultFalse float64        `json:resultFalse`
}
