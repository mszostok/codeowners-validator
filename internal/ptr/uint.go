package ptr

func Uint64Ptr(u uint64) *uint64 {
	c := u
	return &c
}
