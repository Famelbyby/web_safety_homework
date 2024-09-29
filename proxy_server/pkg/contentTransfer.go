package pkg

import (
	"io"
	"log"
	"sync"
)

func Transfer(destination io.WriteCloser, source io.ReadCloser, wg *sync.WaitGroup, wholeBody *string) {
	defer wg.Done()
	defer destination.Close()
	defer source.Close()

	endlessSource := io.TeeReader(source, destination)
	body, err := io.ReadAll(endlessSource)

	if err != nil {
		log.Println(err)
		return
	}

	*wholeBody += string(body)
}
