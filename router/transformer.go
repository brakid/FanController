package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

type Transformer interface {
	TransformData(input Data) []Command
}

type ConfigTransformer struct {
	transformFunctions map[string]TransformationFunctionDescriptor
}

func createTransformer() (Transformer, error) {
	file, _ := ioutil.ReadFile("transformation.json")

	transformFunctions := make(map[string]TransformationFunctionDescriptor)

	err := json.Unmarshal([]byte(file), &transformFunctions)

	if err != nil {
		return nil, err
	}

	return ConfigTransformer{transformFunctions: transformFunctions}, nil
}

func (ct ConfigTransformer) TransformData(input Data) []Command {
	tfd, found := ct.transformFunctions[input.EmitterId]

	if !found {
		log.Printf("No transformation function found: %v", input.EmitterId)
		return make([]Command, 0)
	}

	return applyFunction(tfd, input)
}

func applyFunction(tfd TransformationFunctionDescriptor, input Data) []Command {
	targetIds := tfd.TargetIds

	value, err := evaluate(tfd.FunctionRoot, input.Value)

	commands := make([]Command, len(targetIds))

	if err != nil {
		log.Printf("Error: %v", err)
		return commands[0:0]
	}

	for index, targetId := range targetIds {
		commands[index] = Command{targetId, value}
	}

	return commands
}

func deserialize[N Node | MathNode | ValueNode | ComparatorNode](rawNode json.RawMessage) (*N, error) {
	var root N
	err := json.Unmarshal(rawNode, &root)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling %T: %v", root, err)
	}
	return &root, nil
}

func evaluate(rawRootNode json.RawMessage, input float64) (float64, error) {
	root, err := deserialize[Node](rawRootNode)
	if err != nil {
		return 0.0, err
	}

	switch root.Type {
	case "INPUT_NODE":
		return input, nil
	case "VALUE_NODE":
		return evaluateValueNode(rawRootNode)
	case "MATH_NODE":
		return evaluateMathNode(rawRootNode, input)
	case "COMPARATOR_NODE":
		return evaluateComparatorNode(rawRootNode, input)
	default:
		return 0.0, fmt.Errorf("Invalid Node Type: %v", root.Type)
	}
}

func evaluateValueNode(rawNode json.RawMessage) (float64, error) {
	root, err := deserialize[ValueNode](rawNode)
	if err != nil {
		return 0.0, err
	}

	return root.Value, nil
}

func evaluateMathNode(rawNode json.RawMessage, input float64) (float64, error) {
	root, err := deserialize[MathNode](rawNode)
	if err != nil {
		return 0.0, err
	}

	leftValue, err := evaluate(root.LeftOperand, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating left node: %v", err)
	}

	rightValue, err := evaluate(root.RightOperand, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating right node: %v", err)
	}

	switch root.Operator {
	case ADD:
		return leftValue + rightValue, nil
	case SUBTRACT:
		return leftValue - rightValue, nil
	case MULTIPLY:
		return leftValue * rightValue, nil
	case DIVIDE:
		return leftValue / rightValue, nil
	case MAX:
		return math.Max(leftValue, rightValue), nil
	case MIN:
		return math.Min(leftValue, rightValue), nil
	}

	return 0.0, fmt.Errorf("Invalid Operator: %v", root.Operator)
}

func evaluateComparatorNode(rawNode json.RawMessage, input float64) (float64, error) {
	root, err := deserialize[ComparatorNode](rawNode)
	if err != nil {
		return 0.0, err
	}

	leftValue, err := evaluate(root.LeftOperand, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating left node: %v", err)
	}

	rightValue, err := evaluate(root.RightOperand, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating right node: %v", err)
	}

	resultTrue, err := evaluate(root.ResultTrue, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating true result: %v", err)
	}

	resultFalse, err := evaluate(root.ResultFalse, input)
	if err != nil {
		return 0.0, fmt.Errorf("Error evaluating false result: %v", err)
	}

	var comparatorResult bool
	switch root.Comparator {
	case LOWER:
		comparatorResult = leftValue < rightValue
	case LOWER_EQUAL:
		comparatorResult = leftValue <= rightValue
	case EQUAL:
		comparatorResult = leftValue == rightValue
	case GREATER_EQUAL:
		comparatorResult = leftValue >= rightValue
	case GREATER:
		comparatorResult = leftValue > rightValue
	default:
		return 0.0, fmt.Errorf("Invalid operator")
	}

	if comparatorResult {
		return resultTrue, nil
	} else {
		return resultFalse, nil
	}
}
