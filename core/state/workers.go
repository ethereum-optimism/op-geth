package state

import (
	"errors"
	"runtime"

	"golang.org/x/sync/errgroup"
)

type workerGroup interface {
	Go(func() error)
	SetLimit(int)
	Wait() error
}

func newWorkerGroup() workerGroup {
	if runtime.NumCPU() == 1 {
		return &inlineWorkerGroup{}
	} else {
		var grp errgroup.Group
		return &grp
	}
}

type inlineWorkerGroup struct {
	err error
}

func (i *inlineWorkerGroup) Go(action func() error) {
	i.err = errors.Join(i.err, action())
}

func (i *inlineWorkerGroup) SetLimit(_ int) {
}

func (i *inlineWorkerGroup) Wait() error {
	return i.err
}
