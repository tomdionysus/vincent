package vincent

type Controller func(context *Context) (bool, error)