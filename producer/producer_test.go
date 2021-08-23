package main

import (
	"nsq/producer/mock"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_SendMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockpublisher(ctrl)

	type args struct {
		name    string
		content string
	}
	tests := []struct {
		name string
		args args
	}{
		{"case-1", args{"name_1", "content_1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(m)
			m.EXPECT().Publish("Topic_Example", gomock.Any()).Return(nil)
			h.SendMsg(tt.args.name, tt.args.content)
		})
	}
}
