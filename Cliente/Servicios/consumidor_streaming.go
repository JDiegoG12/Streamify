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

// Usamos sync.Once para asegurar que speaker.Init() se llame una sola vez.
// Llamarlo múltiples veces puede causar problemas con la librería de audio.
var speakerInicializado sync.Once

// Añadimos el parámetro de contexto
func ReproducirCancion(cliente ss.AudioServiceClient, tituloCancion string, ctx context.Context, done chan bool) {
	// Usamos el CONTEXTO que nos pasan. Si se cancela, la llamada gRPC se interrumpirá.
	stream, err := cliente.StreamAudio(ctx, &ss.PeticionDTO{Titulo: tituloCancion})
	if err != nil {
		// Si el error es por cancelación, no es un fallo fatal, es esperado.
		if ctx.Err() == context.Canceled {
			fmt.Println("\nReproducción cancelada por el usuario.")
		} else {
			log.Printf("Error al invocar el streaming: %v", err)
		}
		done <- true
		return
	}

	reader, writer := io.Pipe()

	go recibirFragmentos(stream, writer)
	go decodificarYReproducir(reader, done)
}

func recibirFragmentos(stream ss.AudioService_StreamAudioClient, writer *io.PipeWriter) {
	defer writer.Close()

	for {
		fragmento, err := stream.Recv()
		// Si el error es por cancelación del contexto, stream.Recv() fallará.
		if err == io.EOF {
			return
		}
		if err != nil {
			// No imprimimos el error si el pipe se cierra por una cancelación, es normal.
			if err != io.ErrClosedPipe {
				log.Printf("Error recibiendo fragmento: %v", err)
			}
			return
		}

		if _, err := writer.Write(fragmento.GetData()); err != nil {
			log.Printf("Error escribiendo en pipe: %v", err)
			return
		}
	}
}

func decodificarYReproducir(reader *io.PipeReader, done chan bool) {
	// Siempre notificamos al canal 'done' cuando esta goroutine termina.
	defer func() {
		// reader.Close() también es una buena práctica aquí.
		reader.Close()
		done <- true
	}()

	// mp3.Decode puede leer desde cualquier fuente que implemente io.Reader, como nuestro pipe.
	// Bloqueará hasta que tenga suficientes datos para leer la cabecera del MP3 y determinar el formato.
	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		log.Printf("Error decodificando MP3: %v", err)
		return
	}
	defer streamer.Close()

	// Inicializamos el dispositivo de audio (altavoces) UNA SOLA VEZ.
	speakerInicializado.Do(func() {
		// El buffer de 1/10 de segundo es un buen balance entre latencia y estabilidad.
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})

	// Creamos un canal de control para saber cuándo la reproducción ha finalizado físicamente.
	reproduccionTerminada := make(chan bool)

	// beep.Seq reproduce los streamers en secuencia.
	// beep.Callback es un "streamer" especial que ejecuta una función cuando es su turno.
	// Al ponerlo al final, la función se ejecutará justo cuando 'streamer' termine.
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("\nLa reproducción de audio ha finalizado.")
		reproduccionTerminada <- true
	})))

	// Esperamos hasta que el callback nos notifique que la canción terminó.
	<-reproduccionTerminada
}
