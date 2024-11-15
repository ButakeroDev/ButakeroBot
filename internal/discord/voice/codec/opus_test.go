package codec

import (
	"bytes"
	"context"
	"errors"
	"github.com/Tomas-vilte/GoMusicBot/internal/logging"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, r.err
}

func TestStreamDCAData_ReadsDataCorrectly(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)

	dca := bytes.NewReader([]byte{0x10, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10})
	opusChan := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockLogger.On("Error", "Error EOF o EOF inesperado encontrado durante la transmisión de datos DCA:", mock.AnythingOfType("[]zapcore.Field")).Return()

	err := clientDCA.StreamDCAData(ctx, dca, opusChan, nil)
	if err != nil {
		t.Errorf("StreamDCAData returned an unexpected error: %v", err)
	}

	expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	actual := <-opusChan

	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestStreamDCAData_HandlesEOF(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)
	dca := bytes.NewReader([]byte{})
	opusChan := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())
	mockLogger.On("Error", "Error EOF o EOF inesperado encontrado durante la transmisión de datos DCA:", mock.AnythingOfType("[]zapcore.Field")).Return()
	defer cancel()

	err := clientDCA.StreamDCAData(ctx, dca, opusChan, nil)
	if err != nil {
		t.Errorf("StreamDCAData returned an unexpected error: %v", err)
	}
}

func TestStreamDCAData_CallsPositionCallback(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)
	data := make([]byte, 0)
	for i := 0; i < 100; i++ {
		frame := []byte{0x10, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
		data = append(data, frame...)
	}
	dca := bytes.NewReader(data)
	opusChan := make(chan []byte, 100)
	mockLogger.On("Error", "Error EOF o EOF inesperado encontrado durante la transmisión de datos DCA:", mock.AnythingOfType("[]zapcore.Field")).Return()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callbackChan := make(chan struct{})
	callback := func(position time.Duration) {
		callbackChan <- struct{}{}
	}

	go func() {
		err := clientDCA.StreamDCAData(ctx, dca, opusChan, callback)
		if err != nil {
			t.Errorf("StreamDCAData returned an unexpected error: %v", err)
		}
	}()

	// Consumir los 100 frames enviados al canal
	for i := 0; i < 100; i++ {
		<-opusChan
	}

	// Esperar a que se llame a la función de callback
	select {
	case <-callbackChan:
		// La función de callback fue llamada
	case <-time.After(100 * time.Millisecond):
		t.Error("Position callback was not called")
	}
}
func TestStreamDCAData_HandlesErrorFromReader(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)
	dca := &errorReader{errors.New("test error")}
	opusChan := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockLogger.On("Error", "Error mientras se leia la longitud de DCA", mock.AnythingOfType("[]zapcore.Field")).Return()

	err := clientDCA.StreamDCAData(ctx, dca, opusChan, nil)
	if err == nil {
		t.Error("StreamDCAData should have returned an error")
	}
}

func TestStreamDCAData_HandlesErrorWhileReadingPCM(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)
	dca := bytes.NewReader([]byte{0x04, 0x00, 0x01, 0x02})
	opusChan := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())
	mockLogger.On("Error", "Error mientras se leia PCM de DCA:", mock.AnythingOfType("[]zapcore.Field")).Return()
	defer cancel()

	err := clientDCA.StreamDCAData(ctx, dca, opusChan, nil)
	if err == nil {
		t.Error("StreamDCAData should have returned an error")
	}
}

func TestStreamDCAData_ReturnsNilOnContextCancellation(t *testing.T) {
	mockLogger := new(logging.MockLogger)
	clientDCA := NewDCAStreamerImpl(mockLogger)
	dca := bytes.NewReader([]byte{0x10, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10})
	opusChan := make(chan []byte, 1)

	ctx, cancel := context.WithCancel(context.Background())
	mockLogger.On("Error", "Error EOF o EOF inesperado encontrado durante la transmisión de datos DCA:", mock.AnythingOfType("[]zapcore.Field")).Return()
	cancel() // Cancelar inmediatamente
	defer cancel()

	err := clientDCA.StreamDCAData(ctx, dca, opusChan, nil)
	if err != nil {
		t.Errorf("StreamDCAData returned an unexpected error: %v", err)
	}
}
