package helpers

func Fade(t float64) float64       { return t * t * t * (t*(t*6-15) + 10) }
func Lerp(t, a, b float64) float64 { return a + t*(b-a) }

func Grad(hash int, x, y float64) float64 {
	h := hash & 7

	switch h {
	case 0:
		return x + y // Diagonal Superior Direita
	case 1:
		return -x + y // Diagonal Superior Esquerda
	case 2:
		return x - y // Diagonal Inferior Direita
	case 3:
		return -x - y // Diagonal Inferior Esquerda
	case 4:
		return x // Leste
	case 5:
		return -x // Oeste
	case 6:
		return y // Norte
	case 7:
		return -y // Sul
	default:
		return 0 // Nunca acontece
	}
}
