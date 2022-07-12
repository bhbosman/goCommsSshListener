package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/goprotoextra"
	"golang.org/x/crypto/ssh"
)

func GoRequestChannelHandler(
	name string,
	cancelContext context.Context,
	requestChannel <-chan *ssh.Request,
	ToReactorFunc goprotoextra.ToReactorFunc,
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
						_ = ToReactorFunc(true, request)
					}
				}
			}
			// run channel empty
			for range requestChannel {
			}
		},
	)
}
