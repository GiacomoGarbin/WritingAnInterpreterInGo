package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	env := make(map[string]Object)
	return &Environment{store: env}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, okay := e.store[name]
	return obj, okay
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}