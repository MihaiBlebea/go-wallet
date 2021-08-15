package context

import (
	"errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Cache Service

func init() {
	Cache = New()
}

type Service interface {
	AddContext(userId int64, ctx *Context)
	GetUserContext(userId int64) (*Context, error)
	ResolveStep(resp *tb.Message) string
	HasPendingContext(userId int64) bool
	SkipStep(userId int64) string
	Cancel(userId int64)
}

type service struct {
	cache map[int64]Context
}

func New() Service {
	return &service{
		cache: make(map[int64]Context),
	}
}

func (s *service) AddContext(userId int64, ctx *Context) {
	s.cache[userId] = *ctx
}

func (s *service) GetUserContext(userId int64) (*Context, error) {
	if val, ok := s.cache[userId]; ok {
		return &val, nil
	}

	return &Context{}, errors.New("Could not find user context")
}

func (s *service) HasPendingContext(userId int64) bool {
	if _, ok := s.cache[userId]; ok {
		return true
	}

	return false
}

func (s *service) ResolveStep(resp *tb.Message) string {
	ctx := s.cache[resp.Chat.ID]

	step := ctx.GetCurrentStep()
	step.Process(resp)
	ctx.IncrementStep()

	s.cache[resp.Chat.ID] = ctx

	if ctx.IsComplete() {
		defer s.Cancel(resp.Chat.ID)

		return "Completed"
	}

	return ctx.GetCurrentQuestion()
}

func (s *service) SkipStep(userId int64) string {
	ctx := s.cache[userId]
	ctx.IncrementStep()

	s.cache[userId] = ctx

	return ctx.GetCurrentQuestion()
}

func (s *service) Cancel(userId int64) {
	delete(s.cache, userId)
}

// {
// 	"1234": {
// 		"context": "task-create",
// 		"confirmation": "I added the task with id NEX-1234"
// 		"steps": [
// 			{
// 				"response": "What is the title of the task?"
// 				"payload": "Set up a catch up with Ric"
// 			},
// 			{
// 				"response": "When do you want to schedule this?"
// 				"payload": "20-07-2020"
// 			},
// 			{
// 				"response": "When time?"
// 				"payload": "20-07-2020"
// 			},
// 			{
// 				"response": "Do you want to add any notes?"
// 				"payload": "tell Ric that he must complete his task for this sprint"
// 			},
// 			{
// 				"response": "what is the priority for this task?"
// 				"payload": "1"
// 			},
// 		]
// 	}
// }
