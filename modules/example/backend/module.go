package example

import platformmodule "goravel/app/core/module"

type Module struct{}

func (Module) Name() string { return "example" }

func (Module) Register(app platformmodule.Application) error {
	return app.RegisterService("example.message", "示例模块已注册并启动")
}

func (Module) Boot(_ platformmodule.Application) error { return nil }
