package base

type Option struct {
	Key   string
	Value any
}

type Options []Option

func (o Options) Value(key string) any {
	for i := len(o) - 1; i >= 0; i-- {
		option := o[i]
		if option.Key == key {
			return option.Value
		}
	}

	return nil
}

func extend(dst Options, srcs ...Options) Options {
	size := len(dst)
	for _, src := range srcs {
		size += len(src)
	}

	var extended Options

	if cap(dst) <= size {
		extended = dst
	} else {
		extended = make(Options, 0, size)
		extended = append(extended, dst...)
	}

	for _, src := range srcs {
		extended = append(extended, src...)
	}

	return extended
}
