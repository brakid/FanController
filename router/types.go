package main

import "encoding/json"

type Data struct {
	EmitterId string  `json:"emitterId"`
	Value     float64 `json:"value"`
}

type Command struct {
	TargetId string  `json:"targetId"`
	Value    float64 `json:"value"`
}

const (
	MATH_NODE       string = "MATH_NODE"
	COMPARATOR_NODE string = "COMPARATOR_NODE"
	VALUE_NODE      string = "VALUE_NODE"
	INPUT_NODE      string = "INPUT_NODE"
)

type Node struct {
	Type string `json:"type"`
}

const (
	ADD      string = "ADD"
	SUBTRACT string = "SUBTRACT"
	MULTIPLY string = "MULTIPLY"
	DIVIDE   string = "DIVIDE"
	MAX      string = "MAX"
	MIN      string = "MIN"
)

type MathNode struct {
	Node
	Operator     string          `json:"operator"`
	LeftOperand  json.RawMessage `json:"left"`
	RightOperand json.RawMessage `json:"right"`
}

const (
	LOWER         string = "LOWER"
	LOWER_EQUAL   string = "LOWER_EQUAL"
	EQUAL         string = "EQUAL"
	GREATER_EQUAL string = "GREATER_EQUAL"
	GREATER       string = "GREATER"
)

type ComparatorNode struct {
	Node
	Comparator   string          `json:"comparator"`
	LeftOperand  json.RawMessage `json:"left"`
	RightOperand json.RawMessage `json:"right"`
	ResultTrue   json.RawMessage `json:"resultTrue"`
	ResultFalse  json.RawMessage `json:"resultFalse"`
}

type ValueNode struct {
	Node
	Value float64 `json:"value"`
}

type InputNode struct {
	Node
}

type TransformationFunctionDescriptor struct {
	TargetIds    []string        `json:"targetIds"`
	FunctionRoot json.RawMessage `json:"functionRoot"`
}
