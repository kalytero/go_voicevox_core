# Requirements
- voicevox_core: `0.15.3`
- linux `x64` | `arm64`

# Usage (clean build)
```bash
export BINARY=download-linux-x64
# export BINARY=download-linux-arm64

export LD_LIBRARY_PATH=./voicevox_core

curl -sSfL https://github.com/VOICEVOX/voicevox_core/releases/latest/download/$BINARY -o download && chmod +x download && ./download && go build

./voicevox
```
