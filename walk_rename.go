/*
Example for walk over path (directory) and work with rune.
Edit string type with splitting to rune array.
*/
package main

import (
    "bytes"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "unicode"
)

/*type FileInfo interface {
        Name() string       // base name of the file
        Size() int64        // length in bytes for regular files; system-dependent for others
        Mode() FileMode     // file mode bits
        ModTime() time.Time // modification time
        IsDir() bool        // abbreviation for Mode().IsDir()
        Sys() interface{}   // underlying data source (can return nil)
}*/
var TestFlag bool = false
var exts = []string{".txt", ".zip"}

func NewName(oldName string) string {
    var ooldName = []rune(oldName)
    var buffer bytes.Buffer
    var prev rune = '_'
    var count int = 0
    for _, char := range ooldName {
        if unicode.IsLower(char) && prev == '_' && count < 3 {
            buffer.WriteRune(unicode.ToUpper(char))
            count += 1
        } else if char == '-' {
            buffer.WriteRune('_')
        } else {
            buffer.WriteRune(char)
        }
        prev = char
    }
    return buffer.String()
}

func CheckExt(filename string) bool {
    var e = filepath.Ext(filename)
    for _, item := range exts {
        if item == e {
            return true
        }
    }
    return false
}

func WalkFunc(paths string, info os.FileInfo, err error) error {
    if err != nil {
        fmt.Println(err)
        return nil
    }
    if info.IsDir() {
        return nil
    }
    fmt.Println(filepath.Dir(paths), "\t", info.Name(), "\t", info.Size(), '\t', filepath.Ext(paths))
    return nil
}

func RenamerWalkFunc(oldPath string, info os.FileInfo, err error) error {
    if err != nil {
        fmt.Println(err)
        return nil
    }
    if info.IsDir() {
        return nil
    }
    // You can use Match, but...
    if !CheckExt(oldPath) {
        return nil
    }
    dirname := filepath.Dir(oldPath)
    newName := NewName(info.Name())
    newPath := filepath.Join(dirname, newName)
    if TestFlag == false {
        er := os.Rename(oldPath, newPath)
        if er != nil {
            fmt.Println("Can't rename file: ", oldPath, err)
            return nil
        }
    } else {
        fmt.Println("Rename: ", oldPath, "to", newPath)
    }
    return nil
}

func main() {
    pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    pathPtr := flag.String("path", pwd, "a path wor walk")
    testPtr := flag.Bool("test", TestFlag, "Test, without real rename")
    flag.Parse()
    TestFlag = *testPtr
    filepath.Walk(*pathPtr, RenamerWalkFunc)
}
