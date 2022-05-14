package OKXClient

type Params map[string]string

func ParamsBuilder() Params {
	return make(Params)
}

func (p Params) Set(key, value string) Params {
	p[key] = value
	return p
}
