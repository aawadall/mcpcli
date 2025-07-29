package survey

// MultiSelect is a stub for the real survey.MultiSelect prompt.
type MultiSelect struct {
    Message string
    Options []string
}

// Select provides a list of choices.
type Select struct {
    Message string
    Options []string
    Default interface{}
}

// Confirm asks a yes/no question.
type Confirm struct {
    Message string
    Default bool
}

// Question represents a single survey question.
type Question struct {
    Name     string
    Prompt   interface{}
    Validate interface{}
}

// Required is a placeholder validator for required fields.
var Required = struct{}{}

// Input is a stub for the real survey.Input prompt.
type Input struct {
    Message string
    Default string
}

// AskOne is a function variable to allow overriding in tests.
var AskOne = func(prompt interface{}, response interface{}, opts ...interface{}) error {
    return nil
}

// Ask displays a series of questions.
var Ask = func(qs interface{}, response interface{}, opts ...interface{}) error {
    return nil
}
