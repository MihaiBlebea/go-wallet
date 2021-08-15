package context

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Context struct {
	Label        string
	CurrentStep  int
	Steps        []Step
	Payload      map[string]string
	Confirmation string
}

type Step struct {
	Question func() string
	Process  func(resp *tb.Message)
}

func NewContext(label string) *Context {
	return &Context{
		Label:        label,
		CurrentStep:  0,
		Steps:        make([]Step, 0),
		Payload:      make(map[string]string),
		Confirmation: "All done",
	}
}

func (c *Context) GetCurrentStep() *Step {
	return &c.Steps[c.CurrentStep]
}

func (c *Context) AddStep(step Step) {
	c.Steps = append(c.Steps, step)
}

func (c *Context) GetStep(index int) *Step {
	return &c.Steps[index]
}

func (c *Context) IncrementStep() {
	c.CurrentStep++
}

func (c *Context) GetCurrentQuestion() string {
	return c.GetCurrentStep().Question()
}

func (c *Context) IsComplete() bool {
	return c.CurrentStep == len(c.Steps)
}
