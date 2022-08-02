package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/crypto/ssh"
)

func GoRequestChannelHandler(
	name string,
	cancelContext context.Context,
	requestChannel <-chan *ssh.Request,
	ToReactorFunc rxgo.NextFunc,
	goFunctionCounter GoFunctionCounter.IService,
) error {
	return goFunctionCounter.GoRun(
		fmt.Sprintf("GoRequestChannelHandler.%v", name),
		func() {
		loop:
			for {
				select {
				case <-cancelContext.Done():
					break loop
				case request, ok := <-requestChannel:
					if !ok {
						break loop
					}
					if request != nil {
						ToReactorFunc(request)
					}
				}
			}
			// run channel empty
			for range requestChannel {
			}
		},
	)
}
