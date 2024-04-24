package main

import (
    "fmt"
    "os"
    "bytes"
    "io"
)

func count_num_of_bytes(file []byte) int {
    return len(file)
}

func count_num_of_lines(file []byte) int {
    return bytes.Count(file, []byte{'\n'})
}

func count_num_of_words(file []byte) int {
    return len(bytes.Fields(file))
}

func count_num_of_characters(file []byte) int {
    return len(bytes.Runes(file))
}

type File struct {
    Filename string
    File []byte
}

func read_files(filenames []string) ([]File, error) {
    files := make([]File, 0)

    if len(filenames) == 0 {
        file, err := io.ReadAll(os.Stdin)
        if err != nil {
            return nil, err
        }
        files = append(files, File{"", file})
    }

    for _, filename := range filenames {
        file, err := os.ReadFile(filename)
        if err != nil {
            return nil, err
        }
        files = append(files, File{filename, file})
    }

    return files, nil
}

func wc_print(counts []int, filename string) {
    for _, count := range counts {
        fmt.Printf("%8d", count)
    }
    fmt.Printf(" %v\n", filename)
}

func main() {
    INDICES := map[string]int{
        "-l": 0,
        "-w": 1,
        "-c": 2,
        "-m": 3,
    }

    COUNT_FUNCTIONS := [](func([]byte) int){
        count_num_of_lines,
        count_num_of_words,
        count_num_of_bytes,
        count_num_of_characters,
    }

    flags := make([]bool, 4)
    filenames := make([]string, 0)
    for i, arg := range os.Args[1:] {
        if arg != "-c" && arg != "-l" && arg != "-w" && arg != "-m" {
            filenames = os.Args[i+1:]
            break
        }

        flags[INDICES[arg]] = true
    }

    if flags[0] == false && flags[1] == false && flags[2] == false && flags[3] == false {
        flags[0] = true
        flags[1] = true
        flags[2] = true
    }

    files, err := read_files(filenames)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    valid_functions := make([](func([]byte) int), 0)
    for i, flag := range flags {
        if flag {
            valid_functions = append(valid_functions, COUNT_FUNCTIONS[i])
        }
    }

    totals := make([]int, len(valid_functions))
    for _, file := range files {
        num_counts := make([]int, len(valid_functions))
        for i, function := range valid_functions {
            num_counts[i] = function(file.File)
            totals[i] += num_counts[i]
        }
        wc_print(num_counts, file.Filename)
    }

    if len(files) > 1 {
        wc_print(totals, "total")
    }
}
