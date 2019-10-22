package util

func MustDo(err error) {
	if err != nil {
		panic(err)
	}
}

func MustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleError(err error, h func(...interface{})) {
	if err != nil {
		h(err)
	}
}
