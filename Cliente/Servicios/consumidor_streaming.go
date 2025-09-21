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

func ReproducirCancion(cliente ss.AudioServiceClient, tituloCancion string, done chan bool) {
	stream, err := cliente.StreamAudio(context.Background(), &ss.PeticionDTO{Titulo: tituloCancion})
	if err != nil {
		log.Printf("Error al invocar el streaming: %v", err)
		done <- true // Notificamos que hemos terminado (con error).
		return
	}

	// io.Pipe crea un par de Reader/Writer conectados en memoria.
	// Lo que se escribe en el Writer, puede ser leído desde el Reader.
	// Es perfecto para conectar dos goroutines.
	reader, writer := io.Pipe()

	// Goroutine 1: Recibe los fragmentos de audio del servidor y los escribe en el pipe.
	go recibirFragmentos(stream, writer)

	// Goroutine 2: Lee del pipe, decodifica el MP3 y lo envía a los altavoces.
	go decodificarYReproducir(reader, done)

	fmt.Println("Reproducción iniciada. La aplicación esperará a que la canción termine.")
}

func recibirFragmentos(stream ss.AudioService_StreamAudioClient, writer *io.PipeWriter) {
	defer writer.Close()
	fmt.Println("Recibiendo canción en vivo...")

	// Añadimos un contador para los fragmentos recibidos.
	fragmentoNum := 1

	for {
		fragmento, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\nCanción recibida completa.") // Agregamos un salto de línea para limpiar la salida.
			return
		}
		if err != nil {
			log.Printf("Error recibiendo fragmento: %v", err)
			return
		}

		// Escribimos en el pipe.
		bytesEscritos, err := writer.Write(fragmento.GetData())
		if err != nil {
			log.Printf("Error escribiendo en pipe (probablemente cerrado por el reproductor): %v", err)
			return
		}
		// Imprimimos el log en la consola del cliente.
		// Usamos \r (retorno de carro) para que el mensaje se sobrescriba en la misma línea,
		// creando una bonita animación de carga sin inundar la consola.
		fmt.Printf("\rRecibido fragmento #%d (%d bytes)...", fragmentoNum, bytesEscritos)
		fragmentoNum++

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
