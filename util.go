package sql

func Placeholders(num int) string {
	size := len("?") + len(",")
	buf := make([]byte, num*size-1)
	for i := 0; i < num; i++ {
		buf[i*size] = []byte("?")[0]
		if i != num {
			buf[i*size+1] = []byte(",")[0]
		}
	}
	return string(buf)
}

func IntsToArgs(nums []int64) []interface{} {
	args := make([]interface{}, len(nums))
	for i, n := range nums {
		args[i] = n
	}
	return args
}
