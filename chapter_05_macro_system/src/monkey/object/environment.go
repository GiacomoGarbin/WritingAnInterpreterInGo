package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	env := make(map[string]Object)
	return &Environment{store: env, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, okay := e.store[name]
	if !okay && e.outer != nil {
		obj, okay = e.outer.Get(name)
	}
	return obj, okay
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}