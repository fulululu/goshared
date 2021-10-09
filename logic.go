package goshared

// Ternary if b == true return t else return f
func Ternary(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

// RepeatedlyDo do some operation at least once
// @Param op represent operation function which has 'func() error' signature
// @Param rt represent repeated times
func RepeatedlyDo(op func() error, rt uint) error {
	var count uint = 0
	var err error
	for err = op(); err != nil && count < rt; count++ {
		err = op()
		if err == nil {
			return nil
		}
	}
	return err
}
