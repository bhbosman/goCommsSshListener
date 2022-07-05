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
	// this function is part of the GoFunctionCounter count
	go func(
		cancelContext context.Context,
		requestChannel <-chan *ssh.Request,
		ToReactorFunc goprotoextra.ToReactorFunc,
	) {
		functionName := goFunctionCounter.CreateFunctionName(fmt.Sprintf("GoRequestChannelHandler.%v", name))
		defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
			_ = GoFunctionCounter.Remove(name)
		}(goFunctionCounter, functionName)
		_ = goFunctionCounter.Add(functionName)

		//
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
	}(cancelContext, requestChannel, ToReactorFunc)
	return nil
}
