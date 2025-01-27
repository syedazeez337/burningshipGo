package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
)

func main() {
	const (
		numFrames = 300              // Total frames for the loop
		width     = 1280             // Video width
		height    = 720              // Video height
		maxIter   = 512              // Iterations per pixel
		fps       = 30               // Frames per second
		zoomSpeed = 0.02             // Zoom speed multiplier
	)

	// Create output directory
	os.Mkdir("frames", 0755)

	// Initial viewport parameters (centered on main ship)
	centerX, centerY := -1.75, 0.0
	initialScale := 2.5

	for frame := 0; frame < numFrames; frame++ {
		// Calculate smooth zoom using sine wave for seamless loop
		theta := 2 * math.Pi * float64(frame) / float64(numFrames)
		scale := initialScale * math.Exp(-zoomSpeed*float64(frame)*math.Sin(theta))

		// Calculate dynamic viewport
		xRange := scale * 3.0 // 3:2 aspect ratio
		yRange := scale * 2.0
		
		xMin := centerX - xRange/2
		xMax := centerX + xRange/2
		yMin := centerY - yRange/2
		yMax := centerY + yRange/2

		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// Generate fractal frame
		for px := 0; px < width; px++ {
			for py := 0; py < height; py++ {
				x := xMin + (float64(px)/float64(width))*(xMax-xMin)
				y := yMin + (float64(height-py-1)/float64(height))*(yMax-yMin)

				zx, zy := 0.0, 0.0
				var iteration int

				for iteration = 0; iteration < maxIter && (zx*zx+zy*zy) < 4.0; iteration++ {
					xtmp := zx*zx - zy*zy + x
					zy = 2*math.Abs(zx*zy) + y
					zx = xtmp
				}

				// Animated color based on frame number
				hueShift := float64(frame) / float64(numFrames)
				img.Set(px, py, animatedColor(iteration, maxIter, zx, zy, hueShift))
			}
		}

		// Save frame with sequential numbering
		f, _ := os.Create("frames/frame_" + strconv.Itoa(frame) + ".png")
		png.Encode(f, img)
		f.Close()
	}
}

func animatedColor(iter, maxIter int, zx, zy float64, hueShift float64) color.Color {
	if iter == maxIter {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	// Smooth coloring with animated hue
	modulus := math.Sqrt(zx*zx + zy*zy)
	mu := float64(iter) + 1 - math.Log(math.Log(modulus))/math.Log(2)
	mu /= float64(maxIter)

	// Convert to HSV color space for better color cycling
	hue := math.Mod(0.7+mu+hueShift, 1.0)
	sat := 0.8
	val := math.Pow(mu, 0.3)

	// Convert HSV to RGB
	r, g, b := hsvToRGB(hue, sat, val)
	
	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}

func hsvToRGB(h, s, v float64) (float64, float64, float64) {
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	switch int(i) % 6 {
	case 0: return v, t, p
	case 1: return q, v, p
	case 2: return p, v, t
	case 3: return p, q, v
	case 4: return t, p, v
	default: return v, p, q
	}
}