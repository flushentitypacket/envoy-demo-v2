package main

import (
	"math/rand/v2"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func init() {
	balancer.Register(newRandomLBBuilder())
}

type randomLBBuilder struct{}

func newRandomLBBuilder() balancer.Builder {
	return base.NewBalancerBuilder("random", &randomPickerBuilder{}, base.Config{HealthCheck: true})
}

type randomPickerBuilder struct{}

func (*randomPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &randomPicker{
		subConns: scs,
	}
}

func (b *randomLBBuilder) Name() string {
	return "random"
}

type randomPicker struct {
	subConns []balancer.SubConn
}

func (p *randomPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	subConnsLen := len(p.subConns)
	sc := p.subConns[rand.IntN(subConnsLen)]
	return balancer.PickResult{SubConn: sc}, nil
}
