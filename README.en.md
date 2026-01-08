gobase

English | [中文](README.md)

gobase is a small collection of reusable Go utilities,
organized into focused packages under the utils directory.

Main utility packages:
- arrayutil  : common slice/array helpers (diff, union, subset checks, etc.)
- stringutil: string transformations and random string helpers
- formatutil: formatted output utilities, such as progress bars
- osutil    : OS-related helpers (process name, goroutine ID, default network IP, etc.)
- maputil   : generic map helpers (merge, get with default, keys/values, struct <-> map)

Usage
1. Add the module to your project with:
   go get github.com/ekreke/gobase@latest
2. Import the packages you need in your Go code, for example:
   import "github.com/ekreke/gobase/utils/arrayutil"
   import "github.com/ekreke/gobase/utils/stringutil"
   import "github.com/ekreke/gobase/utils/maputil"

Run tests under utils:
cd utils
go test ./...
