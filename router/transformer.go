package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Transformer interface {
	TransformData(input Data) []Command
}

type ConfigTransformer struct {
	transformFunctions map[string]TransformationFunctionDescriptor
}

func createTransformer() Transformer {
	file, _ := ioutil.ReadFile("transformation.json")

	transformFunctions := make(map[string]TransformationFunctionDescriptor)

	_ = json.Unmarshal([]byte(file), &transformFunctions)

	return ConfigTransformer{transformFunctions: transformFunctions}
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

	value, err := apply(tfd, input.Value)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	commands := make([]Command, len(targetIds))

	for index, targetId := range targetIds {
		commands[index] = Command{targetId, value}
	}

	return commands
}

func apply(tfd TransformationFunctionDescriptor, value float64) (float64, error) {
	switch tfd.Comparator {
	case LOWER:
		if value < tfd.Threshold {
			return tfd.ResultTrue, nil
		} else {
			return tfd.ResultFalse, nil
		}
	case LOWER_EQUAL:
		if value <= tfd.Threshold {
			return tfd.ResultTrue, nil
		} else {
			return tfd.ResultFalse, nil
		}
	case EQUAL:
		if value == tfd.Threshold {
			return tfd.ResultTrue, nil
		} else {
			return tfd.ResultFalse, nil
		}
	case GREATER_EQUAL:
		if value >= tfd.Threshold {
			return tfd.ResultTrue, nil
		} else {
			return tfd.ResultFalse, nil
		}
	case GREATER:
		if value > tfd.Threshold {
			return tfd.ResultTrue, nil
		} else {
			return tfd.ResultFalse, nil
		}
	}
	return 0.0, fmt.Errorf("Invalid operator")
}
