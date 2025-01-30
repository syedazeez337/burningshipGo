package main

import (
    "fmt"
    "image"
    "image/color"
    "image/png"
    "log"
    "math"
    "math/cmplx"
    "os"
    "path/filepath"
    "sync"
)

// Constants for the image dimensions and fractal parameters
const (
    width, height    = 960, 540        // Half HD resolution (for faster rendering)
    maxIterations    = 1000            // Reduced iterations for faster computation
    numFrames        = 600             // 20 seconds at 30 fps
    supersample      = 1               // Disable supersampling for performance
    framesDir        = "output_frames" // Directory to store frames
    maxWorkers       = 16              // Match the number of logical processors
)

func main() {
    // Create the frames directory if it doesn't exist
    err := os.MkdirAll(framesDir, os.ModePerm)
    if err != nil {
        log.Fatalf("Failed to create frames directory: %v", err)
    }

    // Initial zoom window
    xMin, xMax := -2.0, 1.0
    yMin, yMax := -1.5, 1.5

    // Worker pool setup
    var wg sync.WaitGroup
    frameChan := make(chan int, numFrames) // Channel to distribute frames among workers

    // Start worker goroutines
    for w := 0; w < maxWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for frame := range frameChan {
                generateFrame(frame, xMin, xMax, yMin, yMax)
            }
        }()
    }

    // Send frames to the channel
    for frame := 0; frame < numFrames; frame++ {
        frameChan <- frame
    }
    close(frameChan) // Close the channel after all frames are sent

    // Wait for all workers to finish
    wg.Wait()
    log.Println("All frames generated. Use ffmpeg to create a video.")
}

// generateFrame generates a single frame of the Mandelbrot set
func generateFrame(frame int, xMin, xMax, yMin, yMax float64) {
    // Compute the zoom factor for continuous zoom
    zoomFactor := math.Pow(0.99, float64(frame+1)) // Logarithmic zoom

    // Compute the zoom window for this frame
    centerX := (xMin + xMax) / 2
    centerY := (yMin + yMax) / 2
    xRange := (xMax - xMin) * zoomFactor
    yRange := (yMax - yMin) * zoomFactor
    xMin = centerX - xRange/2
    xMax = centerX + xRange/2
    yMin = centerY - yRange/2
    yMax = centerY + yRange/2

    // Slight panning effect (optional)
    centerX += 0.0001 * math.Sin(float64(frame)/10) // Horizontal pan
    centerY += 0.0001 * math.Cos(float64(frame)/10) // Vertical pan
    xMin += centerX - (xMin + xMax) / 2
    xMax += centerX - (xMin + xMax) / 2
    yMin += centerY - (yMin + yMax) / 2
    yMax += centerY - (yMin + yMax) / 2

    // Log progress
    log.Printf("Generating frame %d/%d...", frame+1, numFrames)

    // Create a new image with the specified dimensions
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    // Generate the Mandelbrot set for this frame
    for py := 0; py < height; py++ {
        for px := 0; px < width; px++ {
            r, g, b := uint64(0), uint64(0), uint64(0)
            for sy := 0; sy < supersample; sy++ {
                for sx := 0; sx < supersample; sx++ {
                    // Map sub-pixel position to a point in the complex plane
                    subX := float64(px) + float64(sx)/float64(supersample)
                    subY := float64(py) + float64(sy)/float64(supersample)
                    x0 := subX/float64(width)*(xMax-xMin) + xMin
                    y0 := subY/float64(height)*(yMax-yMin) + yMin
                    c := complex(x0, y0)

                    // Compute the number of iterations for the Mandelbrot set
                    iterations := mandelbrot(c)

                    // Map the iteration count to a color
                    col := getColor(iterations)
                    r += uint64(col.R)
                    g += uint64(col.G)
                    b += uint64(col.B)
                }
            }

            // Average the colors from supersampling
            totalSamples := supersample * supersample
            img.Set(px, py, color.RGBA{
                uint8(r / uint64(totalSamples)),
                uint8(g / uint64(totalSamples)),
                uint8(b / uint64(totalSamples)),
                255,
            })
        }
    }

    // Save the frame as a PNG file in the frames directory
    filename := fmt.Sprintf("frame_%04d.png", frame)
    filePath := filepath.Join(framesDir, filename)
    file, err := os.Create(filePath)
    if err != nil {
        log.Printf("Failed to create file '%s': %v", filePath, err)
        return
    }
    defer file.Close()

    err = png.Encode(file, img)
    if err != nil {
        log.Printf("Failed to encode image '%s': %v", filePath, err)
        return
    }

    log.Printf("Frame %d/%d saved as '%s'", frame+1, numFrames, filePath)
}

// mandelbrot computes the number of iterations for a given complex number
func mandelbrot(c complex128) int {
    var z complex128
    for i := 0; i < maxIterations; i++ {
        z = z*z + c
        if cmplx.Abs(z) > 2 {
            return i
        }
    }
    return maxIterations
}

// getColor maps the iteration count to a vibrant color
func getColor(iterations int) color.RGBA {
    if iterations == maxIterations {
        // Points inside the Mandelbrot set are white
        return color.RGBA{255, 255, 255, 255}
    }

    // Simple color mapping without dynamic cycling
    t := float64(iterations) / float64(maxIterations)
    r := uint8(255 * (1 - t))          // Red decreases as t increases
    g := uint8(255 * t * t)            // Green increases quadratically
    b := uint8(255 * t * t * t)        // Blue increases cubically

    return color.RGBA{r, g, b, 255}
}