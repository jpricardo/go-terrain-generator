# go-terrain-generator

A procedural terrain generator written in Go. This project generates 2D terrain maps using Perlin Noise, applying elevation-based materials and exporting the result as PNG images.

> **Note:** This project is a **learning exercise** and is not intended for production use. It serves as an exploration of Go, image manipulation, and procedural generation algorithms.

## 🌟 Features

* **Procedural Generation:** Uses **Perlin Noise** with fractal Brownian motion (multiple octaves) to create organic-looking heightmaps.
* **Chunk System:** Generates terrain in chunks to manage processing of larger maps.
* **Material System:** Automatically assigns terrain types based on elevation:
    * Water
    * Sand
    * Grass
    * Stone
    * Snow
* **Image Export:** Renders the final terrain to PNG format.
* **Customizable:** Adjustable seed, smoothness, and map size.

## 🛠️ Installation

1.  **Prerequisites:** Ensure you have Go installed (this project uses Go 1.25.4).
2.  **Clone the repository:**
    ```bash
    git clone https://github.com/jpricardo/go-terrain-generator.git
    cd go-terrain-generator
    ```
3.  **Download dependencies:**
    ```bash
    go mod download
    ```

## 🚀 Usage

To generate a new terrain map, run the `main.go` file inside the CLI directory:

```bash
go run cmd/cli/main.go
```

### Output
The generated images will be saved in the `cmd/output/` directory.

* `chunked_terrain.png`: The final combined terrain map.
* `chunk_x_y.png`: Individual chunk segments.

## ⚙️ Configuration
You can tweak the generation parameters in `cmd/cli/main.go`:

```go
func main() {
    size := 512           // The size of the map
    smoothness := 0.75    // Controls the "zoom" of the noise (higher = smoother/zoomed out)
    seed := uint64(...)   // The seed for the random number generator
    // ...
}
```

You can also adjust the elevation thresholds for materials in `cmd/data/terrain.go`:

```go
const (
    WaterLevel  = uint8(0)
    GroundLevel = uint8(96)
    SnowLevel   = uint8(192)
    // ...
)
```

## 📂 Project Structure
* **`cmd/cli`**: Entry points for the application (Main logic and generator orchestration).
* **`cmd/data`**: Core data structures and algorithms.
* `noise.go`: Implementation of Perlin Noise and White Noise.
* `terrain.go`: Logic for Terrain points, Chunks, and Material application.
* `material.go`: Definitions of colors and material types.
* `texture.go`: Image merging and filter application logic.


* **`cmd/helpers`**: Utility functions for geometry (`Lerp`, `Fade`), hashing, and file I/O.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.