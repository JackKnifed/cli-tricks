package main

import (
	"flag"
	"io"
	"bufio"
	"os"
	"encoding/json"
	"log"
)

func PrettyPrint(in io.Reader, out io.Writer, indent string) (err error) {
	var stuff interface{}
	var writeBuf []byte
	decoder := json.NewDecoder(in)
	for decoder.More() {
		err = decoder.Decode(&stuff)
		if err != nil {
			return err
		}
		writeBuf, err = json.MarshalIndent(stuff, "", indent)
		if err != nil {
			return err
		}
		out.Write(writeBuf)
	}
	if err == io.EOF {
		return nil
	} else {
		return err
	}
}

var indentVal = flag.String("indent", "  ", "indentation to use")

func main() {
	flag.Parse()
	if err := PrettyPrint(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout), *indentVal); err != nil {
		log.Fatal(err)
	}
}
