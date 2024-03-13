# pixelate

Image Processing Service

## Overview

This service provodive the following functionalities:

1. Convert image files from PNG to JPEG.
2. Resize images according to specified dimensions.
3. Compress images to reduce file size while maintaining reasonable quality.

## Prerequisites

### Installation of ffmpeg

```bash
# On Linux (Ubuntu/Debian)
sudo apt-get install ffmpeg

# On macOS (via Homebrew)
brew install ffmpeg

# On Windows
# Download ffmpeg from [official website](https://ffmpeg.org/download.html)
# Add ffmpeg to the system PATH
```

## Endpoints

### Convert

- Description: Convert image files from PNG to JPEG
- Path: `/convert`
- Method: `POST`
- Request Body:
  - `image`: The file to be converted. (Multipart request body)
- Response: The converted file in .JPG

#### Example Usage

```bash
curl -X POST \
  -F "image=example.jpg" \
  https://{host}:{port}/convert
```

### Resize

- Description: Resize images according to specified dimensions
- Path: `/resize`
- Method: `POST`
- Request Body:
  - `image`: The file to be converted. (Multipart request body)
  - `scale`: specified dimensions image
- Response: The file with specified dimensions image

#### Example Usage

```bash
curl -X POST \
  -F "image=example.jpg" \
  -F "scale=640:640" \
  https://{host}:{port}/resize
```

### Compress

- Description: Compress images to reduce file size while maintaining reasonable quality
- Path: `/compress`
- Method: `POST`
- Request Body:
  - `image`: The file to be converted. (Multipart request body)
- Response: The reduced file

#### Example Usage

```bash
curl -X POST \
  -F "image=example.jpg" \
  https://{host}:{port}/compress
