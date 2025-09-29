// Archivo: Cliente/Servicios/streaming_servicios.go
package Servicios

import (
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var speakerInicializado sync.Once

// La firma de la función ya es correcta, pasa el contexto.
func ReproducirCancion(cliente ss.AudioServiceClient, tituloCancion string, ctx context.Context, done chan bool) {
	stream, err := cliente.StreamAudio(ctx, &ss.PeticionDTO{Titulo: tituloCancion})
	if err != nil {
		if ctx.Err() != context.Canceled {
			log.Printf("Error al invocar el streaming: %v", err)
		}
		done <- true
		return
	}

	reader, writer := io.Pipe()

	// La goroutine de recepción de red se mantiene igual (versión silenciosa).
	go recibirFragmentos(stream, writer)

	// CAMBIO CLAVE: Pasamos el 'ctx' a la goroutine de decodificación y reproducción.
	go decodificarYReproducir(reader, ctx, done)
}

// Versión final y silenciosa de recibirFragmentos.
func recibirFragmentos(stream ss.AudioService_StreamAudioClient, writer *io.PipeWriter) {
	defer writer.Close()
	for {
		fragmento, err := stream.Recv()
		if err != nil {
			// Cualquier error (EOF, cancelado, red) detiene esta goroutine silenciosamente.
			return
		}

		if _, err := writer.Write(fragmento.GetData()); err != nil {
			// Si el pipe se cierra, también terminamos.
			return
		}
	}
}

// CAMBIO CLAVE: Esta función ahora acepta el 'ctx' y usa un 'select' para evitar el bloqueo.
func decodificarYReproducir(reader *io.PipeReader, ctx context.Context, done chan bool) {
	// El defer asegura que siempre se notifique al canal 'done' cuando esta goroutine termine.
	defer func() {
		reader.Close()
		done <- true
	}()

	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		// Este error es normal si se cancela antes de recibir suficientes datos de audio.
		return
	}
	defer streamer.Close()

	speakerInicializado.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})

	// Canal para saber cuándo la reproducción de beep ha terminado de forma natural.
	reproduccionTerminada := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// Cuando la canción termina, cerramos el canal para señalarlo.
		close(reproduccionTerminada)
	})))

	// 'select' espera por el primer evento que ocurra:
	// 1. La canción termina por sí sola.
	// 2. El usuario cancela la operación.
	select {
	case <-reproduccionTerminada:
		// La canción terminó naturalmente. Imprimimos el mensaje de finalización.
		fmt.Println("\nLa reproducción de audio ha finalizado.")
	case <-ctx.Done():
		// El contexto fue cancelado por el usuario. No imprimimos nada y simplemente
		// salimos de la función. El 'defer' se encargará de la limpieza.
		return
	}
}
