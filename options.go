package tqla

type options struct {
	placeholder Placeholder
}

type Option interface {
	Apply(*options) error
}

type funcOption struct {
	f func(*options) error
}

func (flo *funcOption) Apply(con *options) error {
	return flo.f(con)
}

func newFuncNodeOption(f func(*options) error) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithPlaceHolder(p Placeholder) Option {
	return newFuncNodeOption(func(o *options) error {
		o.placeholder = p
		return nil
	})
}

func defaultOptions() *options {
	return &options{
		placeholder: Question,
	}
}
