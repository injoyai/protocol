package fins

type Pkg struct {
}

func (this *Pkg) Bytes() []byte {
	data := []byte(nil)
	data = append(data, '@')

	data = append(data, '*', '\n')
	return data
}
