package Servicios

import (
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// El resto de este archivo no necesita cambios.
func ReproducirCancion(cliente ss.AudioServiceClient, tituloCancion string) {
	stream, err := cliente.StreamAudio(context.Background(), &ss.PeticionDTO{Titulo: tituloCancion})
	if err != nil {
		log.Printf("Error al invocar el streaming: %v. Asegúrate de que el archivo '%s.mp3' existe en el ServidorStreaming.", err, tituloCancion)
		return
	}

	reader, writer := io.Pipe()

	go recibirFragmentos(stream, writer)
	go decodificarYReproducir(reader)

	fmt.Println("Reproducción iniciada. Presiona Enter para volver al menú anterior cuando la canción termine.")
}

func recibirFragmentos(stream ss.AudioService_StreamAudioClient, writer *io.PipeWriter) {
	defer writer.Close()
	fmt.Println("Recibiendo canción en vivo...")

	for {
		fragmento, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Canción recibida completa.")
			return
		}
		if err != nil {
			log.Fatalf("Error recibiendo fragmento: %v", err)
		}

		_, err = writer.Write(fragmento.GetData())
		if err != nil {
			log.Printf("Error escribiendo en pipe: %v", err)
			return
		}
	}
}

func decodificarYReproducir(reader *io.PipeReader) {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("¡¡¡PÁNICO ATRAPADO EN LA GOROUTINE DE AUDIO!!! Error: %v", r)
		}
	}()

	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		log.Fatalf("Error decodificando MP3: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))
	speaker.Play(streamer)
}
