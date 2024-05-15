package toolbox

import "math/rand"

func RandomFill(buf []byte, n int) {
	for i := 0; i < min(n, len(buf)); i++ {
		if rand.Int31n(2) == 0 {
			buf[i] = byte(65 + rand.Int31n(26))
		} else {
			buf[i] = byte(97 + rand.Int31n(26))
		}
	}
}
