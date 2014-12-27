/*
Work with conversion between rune string and modify it.
*/
package main

import (
    "bytes"
    "fmt"
    "unicode"
    "unicode/utf8" // For NewName2
    "bufio" // For NewName3
)

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
    // AT! You don't want to use string([]rune)!
    return buffer.String()
}

func NewName2(oldName string) string {
    rune_count := len(oldName)
    var prev rune = '_'
    var count int = 0
    var buffer bytes.Buffer
    for i := 0; i < rune_count; {
        if i < rune_count {
            r, w := utf8.DecodeRuneInString(oldName[i:])
            i += w
            if unicode.IsLower(r) && prev == '_' && count < 3 {
                buffer.WriteRune(unicode.ToUpper(r))
                count += 1
            } else if r == '-' {
                buffer.WriteRune('_')
            } else {
                buffer.WriteRune(r)
            }
            prev = r
        }
    }
    return buffer.String()
}

func NewName3(oldName string) string {
    rune_count := utf8.RuneCountInString(oldName)
    var prev rune = '_'
    var count int = 0
    var buffer bytes.Buffer
    input := bufio.NewReader(bytes.NewBufferString(oldName));
    for i := 1; i <= rune_count; i++ {
        r, _, _ := input.ReadRune();
        if unicode.IsLower(r) && prev == '_' && count < 3 {
            buffer.WriteRune(unicode.ToUpper(r))
            count += 1
        } else if r == '-' {
            buffer.WriteRune('_')
        } else {
            buffer.WriteRune(r)
        }
        prev = r
    }
    return buffer.String()
}

func main() {
}
