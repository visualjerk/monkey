package object

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outterEnv *Environment) *Environment {
	return &Environment{
		store:     make(map[string]Object),
		outterEnv: outterEnv,
	}
}

type Environment struct {
	store     map[string]Object
	outterEnv *Environment
}

func (env *Environment) Get(key string) (Object, bool) {
	object, ok := env.store[key]

	if !ok && env.outterEnv != nil {
		return env.outterEnv.Get(key)
	}

	return object, ok
}

func (env *Environment) Set(key string, value Object) Object {
	env.store[key] = value
	return value
}
