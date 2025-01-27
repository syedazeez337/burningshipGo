package dummy

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	const (
		width   = 1600
		height  = 1200
		maxIter = 512
	)

	// Viewport coordinates (adjust these to explore different regions)
	xMin, xMax := -2.0, 1.5
	yMin, yMax := -1.8, 0.8

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for px := 0; px < width; px++ {
		for py := 0; py < height; py++ {
			// Map pixel coordinates to complex plane
			x := xMin + (float64(px)/float64(width))*(xMax-xMin)
			y := yMin + (float64(height-py-1)/float64(height))*(yMax-yMin)

			zx, zy := 0.0, 0.0 // Initialize Z = 0
			var iteration int

			for iteration = 0; iteration < maxIter && (zx*zx+zy*zy) < 4.0; iteration++ {
				xtmp := zx*zx - zy*zy + x
				zy = 2*math.Abs(zx*zy) + y
				zx = xtmp
			}

			// Capture final z values for smooth coloring
			finalZx, finalZy := zx, zy
			img.Set(px, py, fractalColor(iteration, maxIter, finalZx, finalZy))
		}
	}

	// Save image
	f, _ := os.Create("burning_ship.png")
	defer f.Close()
	png.Encode(f, img)
}

// Updated color function with proper parameters
/*
func fractalColor(iter, maxIter int, zx, zy float64) color.Color {
	if iter == maxIter {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Inside points: black
	}

	// Smooth coloring algorithm
	modulus := math.Sqrt(zx*zx + zy*zy)
	mu := float64(iter) + 1 - math.Log(math.Log(modulus))/math.Log(2)
	mu /= float64(maxIter)

	// Color palette parameters (adjust these for different effects)
	r := uint8(255 * math.Pow(math.Sin(mu*math.Pi*0.7), 2))
	g := uint8(255 * math.Pow(math.Sin(mu*math.Pi*1.3), 2))
	b := uint8(255 * math.Pow(math.Sin(mu*math.Pi*2.1), 2))

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

/*
func fractalColor(iter, maxIter int, zx, zy float64) color.Color {
	if iter == maxIter {
		// Inside points: black
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	// Smooth coloring
	modulus := math.Sqrt(zx*zx + zy*zy)
	mu := float64(iter) + 1 - math.Log(math.Log(modulus))/math.Log(2)
	mu /= float64(maxIter)

	// Gradient-based color mapping
	r := uint8(255 * math.Pow(math.Sin(mu*math.Pi*1.0), 2)) // Red channel
	g := uint8(255 * math.Pow(math.Sin(mu*math.Pi*2.0), 2)) // Green channel
	b := uint8(255 * math.Pow(math.Sin(mu*math.Pi*3.0), 2)) // Blue channel

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
*/


func fractalColor(iter, maxIter int, zx, zy float64) color.Color {
	if iter == maxIter {
		// Inside points: black
		return color.Black
	}

	// Smooth escape time
	modulus := math.Sqrt(zx*zx + zy*zy)
	mu := float64(iter) + 1 - math.Log(math.Log(modulus))/math.Log(2)

	// Normalize to [0, 1]
	normalized := mu / float64(maxIter)

	// Introduce more variation in colors
	h := (normalized * 360) + (float64(iter)*0.1) // Hue varies with normalized value and iterations
	s := 0.8 + 0.2*math.Sin(float64(iter))        // Slightly vary saturation
	v := 0.9                                      // Brightness is fixed for vividness

	// Convert HSV to RGB
	r, g, b := hsvToRgb(h, s, v)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// hsvToRgb converts HSV to RGB (unchanged)
func hsvToRgb(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

