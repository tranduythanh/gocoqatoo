package coq

type InputOutput struct {
	Input  *Input
	Output *Output
}

func NewInputOutput(input *Input, output *Output) *InputOutput {
	return &InputOutput{
		Input:  input,
		Output: output,
	}
}

func (io *InputOutput) GetInput() *Input {
	return io.Input
}

func (io *InputOutput) GetOutput() *Output {
	return io.Output
}
