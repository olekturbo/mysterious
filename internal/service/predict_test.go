package service

import (
	"context"
	"errors"
	"testing"
)

type mockPredictor struct {
	predictFn func(ctx context.Context, text string) (string, error)
}

func (m *mockPredictor) Predict(ctx context.Context, text string) (string, error) {
	return m.predictFn(ctx, text)
}

func TestPredict_Predict(t *testing.T) {
	type fields struct {
		predictor predictor
	}
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "err",
			fields: fields{
				predictor: &mockPredictor{
					predictFn: func(ctx context.Context, text string) (string, error) {
						return "", errors.New("err")
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				predictor: &mockPredictor{
					predictFn: func(ctx context.Context, text string) (string, error) {
						return text, nil
					},
				},
			},
			args: args{
				text: "ok",
			},
			want:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Predict{
				predictor: tt.fields.predictor,
			}
			got, err := p.Predict(context.Background(), tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Predict.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Predict.Predict() = %v, want %v", got, tt.want)
			}
		})
	}
}
