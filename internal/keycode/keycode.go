package keycode

// ASCII keycodes
// ref: https://www.ascii-code.com
var (
	J         = [3]byte{106, 0, 0}
	K         = [3]byte{107, 0, 0}
	UP        = [3]byte{27, 91, 65}
	DOWN      = [3]byte{27, 91, 66}
	ENTER     = [3]byte{13, 0, 0}
	ESCAPE    = [3]byte{27, 0, 0}
	CONTROL_C = [3]byte{3, 0, 0}
	D         = [3]byte{100, 0, 0}
	Y         = [3]byte{121, 0, 0}
	Q         = [3]byte{113, 0, 0}
)

func Matches(input []byte, keycodes ...[3]byte) bool {
	for _, keycode := range keycodes {
		if input[0] == keycode[0] && input[1] == keycode[1] && input[2] == keycode[2] {
			return true
		}
	}

	return false
}
