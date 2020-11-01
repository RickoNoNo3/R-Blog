package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"./cmd"
	"./data"
)

var out string // 最终的stdout

func main() {
	defer data.CloseDB()
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "draw":
			if len(os.Args) >= 3 {
				id, err := strconv.Atoi(os.Args[2])
				if err == nil {
					out = cmd.Draw(id)
				}
			}
		case "drawCore":
			readFromFile := false
			if len(os.Args) >= 3 {
				inputFile, err := os.Open(os.Args[2])
				if err == nil {
					readFromFile = true
					out = cmd.DrawCore(inputFile)
					_ = inputFile.Close()
				}
			}
			if !readFromFile {
				out = cmd.DrawCore(os.Stdin)
			}
		case "new":
			if len(os.Args) >= 3 {
				var (
					err        error
					entityType = 0
					dirId      = 0
				)
				entityType, err = strconv.Atoi(os.Args[2])
				if err == nil {
					if len(os.Args) >= 4 {
						dirId, err = strconv.Atoi(os.Args[3])
					}
					if err == nil {
						out = cmd.New(entityType, dirId, os.Stdin)
					}
				}
			}
		case "drag":
			if len(os.Args) >= 3 {
				if id, err := strconv.Atoi(os.Args[2]); err == nil {
					out = cmd.Drag(id)
				}
			}
		case "edit":
			if len(os.Args) >= 4 {
				var (
					err        error
					entityType = 0
					entityId   = 0
				)
				entityType, err = strconv.Atoi(os.Args[2])
				if err == nil {
					entityId, err = strconv.Atoi(os.Args[3])
					if err == nil {
						out = cmd.Edit(
							data.Entity{
								Type: entityType,
								Id:   entityId,
							}, os.Stdin,
						)
					}
				}
			}
		case "move":
			dirId := 0
			if len(os.Args) >= 3 {
				dirId, _ = strconv.Atoi(os.Args[2])
			}
			out = cmd.Move(dirId, os.Stdin)
		case "remove":
			out = cmd.Remove(os.Stdin)
		case "read":
			if len(os.Args) >= 4 {
				if entityType, err := strconv.Atoi(os.Args[2]); err == nil {
					if id, err := strconv.Atoi(os.Args[3]); err == nil {
						out = cmd.Read(
							data.Entity{
								Type: entityType,
								Id:   id,
							},
						)
					}
				}
			}
		case "link":
			if len(os.Args) >= 4 {
				if entityType, err := strconv.Atoi(os.Args[2]); err == nil {
					if id, err := strconv.Atoi(os.Args[3]); err == nil {
						out = cmd.Link(
							data.Entity{
								Type: entityType,
								Id:   id,
							},
						)
					}
				}
			}
		}
	}
	if out == "" {
		file, err := os.Open("help.txt")
		if err == nil {
			_, _ = io.Copy(os.Stdout, file)
			_ = file.Close()
		}
	} else {
		fmt.Println(out)
	}
}
