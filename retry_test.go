package gpt3

import (
	"errors"
	"testing"
	"time"
)

func Test_retry(t *testing.T) {
	type args struct {
		fn      RetryHandle
		retries int
		delay   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{func() error {
			t.Log("time", time.Now())
			return errors.New("")
		}, 0, time.Second * 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.fn, tt.args.retries, tt.args.delay); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
