package object

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

type Environment struct {
	store map[string]Object
}

func (env *Environment) Get(identifier string) (Object, bool) {
	object, ok := env.store[identifier]
	return object, ok
}

func (env *Environment) Set(identifier string, value Object) Object {
	env.store[identifier] = value
	return value
}
