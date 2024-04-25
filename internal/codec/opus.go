package codec

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	frameLength = time.Duration(20) * time.Millisecond
)

func StreamDCAData(ctx context.Context, dca io.Reader, opusChan chan<- []byte, positionCallback func(position time.Duration)) error {
	var opuslen int16
	framesSent := 0

	for {

		err := binary.Read(dca, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("while reading length from DCA: %w", err)
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(dca, binary.LittleEndian, &inBuf)

		if err != nil {
			return fmt.Errorf("while reading PCM from DCA: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case opusChan <- inBuf:
			framesSent += 1
			go func() {
				if positionCallback != nil && framesSent%50 == 0 {
					positionCallback(time.Duration(framesSent) * frameLength)
				}
			}()
			continue
		}
	}
}
