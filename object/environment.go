package object

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

type Environment struct {
	store map[string]Object
}

func (env *Environment) Get(key string) (Object, bool) {
	object, ok := env.store[key]
	return object, ok
}

func (env *Environment) Set(key string, value Object) Object {
	env.store[key] = value
	return value
}

func (env *Environment) Merge(otherEnv *Environment) {
	for key, value := range otherEnv.store {
		env.store[key] = value
	}
}
