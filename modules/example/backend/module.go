package example

import platformmodule "goravel/app/core/module"

type Module struct{}

func (Module) Name() string { return "example" }

func (Module) Register(_ platformmodule.Application) error { return nil }

func (Module) Boot(_ platformmodule.Application) error { return nil }
