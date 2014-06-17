/*
Convert string "key1:val1:key2:val2"
to map object.

author Yuri Astrov (yuriastrov@gmail.com)
*/
package main

import (
    "fmt"
    "strings"
    "encoding/json"
)

type MyError struct {
    What string
}

func (e MyError) Error() string {
    return fmt.Sprintf("%v", e.What)
}

// Naive realization
func StrToMap(str string) (map[string]string, error) {
    if !strings.Contains(str, ":") {
        return nil, MyError{
            "No delimeter ':' in string!",
        }
    }
    result := strings.Split(str, ":")
    dct := make(map[string]string, len(result))
    for i := 0; i < len(result)-1; i += 2 {
        dct[result[i]] = result[i+1]
    }
    return dct, nil
}

func StrToMap2(str string) (map[string]string, error) {
    if !strings.Contains(str, ":") {
        return nil, MyError{
            "No delimeter ':' in string!",
        }
    }
    result := strings.Split(str, ":")
    //dct := map[string]string{}
    dct := make(map[string]string, len(result))
    var key, val string
    for k, v := range result {
        if k % 2 == 0 {
            key = v
        } else {
            val = v
            dct[key] = val
        }
    }
    return dct, nil
}

/* Best solution in code: Cycle counter version. :)
But may have unnecessary step (operation): dct[result[k-1] ] = v
*/
func StrToMap3(str string) (map[string]string, error) {
    if !strings.Contains(str, ":") {
        return nil, MyError{
            "No delimeter ':' in string!",
        }
    }
    result := strings.Split(str, ":")
    //dct := map[string]string{}
    dct := make(map[string]string, len(result))
    for k, v := range result {
        if k % 2 != 0 {
            dct[result[k-1] ] = v
        }
    }
    return dct, nil
}

/*Cycle counter version.*/
func StrToRichMap(str string) (map[string][]string, error) {
    if !strings.Contains(str, ":") {
        return nil, MyError{
            "No delimeter ':' in string!",
        }
    }
    result := strings.Split(str, ":")
    dct := make(map[string][]string, len(result))
    var key string
    for k, v := range result {
        if k % 2 == 0 {
            key = v
        } else {
            dct[key] = append(dct[key], v)
        }
    }
    return dct, nil
}

func PrintMap(mp map[string]string) {
    out, err := json.Marshal(mp)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(out))
}

func PrintMap2(mp map[string]string) {
    fmt.Printf("Map:\n")
    for k, v := range mp {
        fmt.Printf("%s -> %s\n", k, v)
    }
}

func PrintRichMap(mp map[string][]string) {
    for key, values := range mp {
        fmt.Printf("Key -> %s\n", key)
        fmt.Printf("Val -> %s\n", strings.Join(values, ", "))
    }
}

func main() {
    var str string = "key1:val1:key2:val2:key1:val3"
    fmt.Println("---------")
    fmt.Println("Output as JSON:")
    dct1, err := StrToMap(str)
    if err != nil {
        fmt.Println(err)
    }
    PrintMap(dct1)
    fmt.Println("---------")
    fmt.Println("Iteration:")
    dct2, err := StrToMap2(str)
    if err != nil {
        fmt.Println(err)
    }
    PrintMap2(dct2)
    fmt.Println("---------")
    dct3, err := StrToMap3(str)
    if err != nil {
        fmt.Println(err)
    }
    PrintMap2(dct3)
    fmt.Println("---------")
    fmt.Println("RichMap example")
    dct4, err := StrToRichMap(str)
    if err != nil {
        fmt.Println(err)
    }
    PrintRichMap(dct4)
}
