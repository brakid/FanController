package main

type Transformer interface {
	TransformData(input Data) []Command
}

type DummyTransformer struct{}

func (dt DummyTransformer) TransformData(input Data) []Command {
	commands := make([]Command, 1)
	commands[0] = Command{"DummyTarget", "DummyCommand"}

	return commands
}
