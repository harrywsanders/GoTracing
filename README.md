# GoTracing

**GoTracing** is a simple ray-tracing program written in Go that brings your digital dreams to life! With GoTracing, you can create stunning computer-generated images by simulating the behavior of light rays as they interact with virtual objects. It's like painting with photons!

## Features

- 🌞 Powerful ray-tracing algorithm that brings your scenes to life with realistic lighting and shading.
- 🖼️ Flexible scene configuration with support for custom objects, materials, lights, and camera settings.
- 🚀 Concurrent processing to leverage the full power of your multi-core CPU and render images faster.
- 🌈 Create beautiful scenes with reflective surfaces, transparent materials, and stunning visual effects.
- 🎨 Save your rendered images in popular formats such as PNG or JPEG to showcase your artistic creations.

## Getting Started

### Prerequisites

To run GoTracing, you need to have Go installed on your machine. Check out the official Go installation guide at [golang.org](https://golang.org/doc/install) for instructions.

### Installation

1. Clone the GoTracing repository to your local machine:

```shell
git clone https://github.com/your-username/gotracing.git
```

2. Navigate to the project directory:

```shell
cd gotracing
```

3. Build the project:

```shell
go build
```

## Usage
To render a scene, execute the following command:
```shell
./gotracing render -scene <path_to_scene_file> -output <path_to_output_image>
```
Replace <path_to_scene_file> with the path to the scene file you want to render, and <path_to_output_image> with the desired location and filename for the rendered image.

For example:
```shell
./gotracing render -scene scenes/scene1.json -output images/output.png
```

As always, feel free to contribute, and a huge shoutout to all the computer graphics researchers and artists who inspire us with their work and make super cool things out of this technology!

