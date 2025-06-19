package service

import "context"

type predictor interface {
	Predict(ctx context.Context, text string) (string, error)
}

type Predict struct {
	predictor predictor
}

func NewPredict(predictor predictor) *Predict {
	return &Predict{
		predictor: predictor,
	}
}

func (p *Predict) Predict(ctx context.Context, text string) (string, error) {
	return p.predictor.Predict(ctx, text)
}
