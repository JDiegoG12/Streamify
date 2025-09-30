// El paquete Servicios encapsula la lógica de comunicación del cliente con los
// microservicios remotos. Este archivo se especializa en manejar la complejidad
// del streaming de audio: recibir fragmentos de datos, decodificarlos en tiempo real
// y reproducirlos a través del hardware de sonido del sistema.
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

// speakerInicializado utiliza sync.Once para garantizar que el hardware de audio
// del sistema sea inicializado una única vez durante el ciclo de vida de la aplicación.
// Esto previene errores o comportamientos inesperados que podrían surgir al intentar
// inicializar el dispositivo de audio múltiples veces.
var speakerInicializado sync.Once

// ReproducirCancion es la función principal que orquesta el proceso de streaming de audio.
// Implementa una arquitectura de pipeline concurrente utilizando dos goroutines que se
// comunican a través de un 'io.Pipe':
//
//  1. Goroutine de Red ('recibirFragmentos'): Se conecta al servidor gRPC, recibe los
//     fragmentos de audio y los escribe en el 'PipeWriter'.
//
//  2. Goroutine de Audio ('decodificarYReproducir'): Lee los datos del 'PipeReader',
//     los decodifica de formato MP3 y los envía al altavoz para su reproducción.
//
// Este enfoque desacopla la E/S de red de la decodificación de audio, permitiendo que
// la reproducción sea fluida incluso si hay variaciones en la velocidad de la red.
//
// Parámetros:
//   - cliente: El cliente gRPC para el servicio de streaming.
//   - tituloCancion: El nombre de la canción a solicitar al servidor.
//   - ctx: Un contexto que permite cancelar la operación de streaming desde fuera
//     (por ejemplo, cuando el usuario detiene la reproducción).
//   - done: Un canal que se utiliza para señalar que ambas goroutines han finalizado
//     su ejecución y limpiado sus recursos.
func ReproducirCancion(cliente ss.AudioServiceClient, tituloCancion string, ctx context.Context, done chan bool) {
	// Inicia la llamada de streaming RPC al servidor.
	stream, err := cliente.StreamAudio(ctx, &ss.PeticionDTO{Titulo: tituloCancion})
	if err != nil {
		// Solo registra el error si no fue causado por una cancelación explícita del usuario.
		if ctx.Err() != context.Canceled {
			log.Printf("Error al invocar el streaming: %v", err)
		}
		done <- true // Señala la finalización inmediata si la conexión falla.
		return
	}

	// Crea un pipe en memoria. Los datos escritos en 'writer' pueden ser leídos desde 'reader'.
	reader, writer := io.Pipe()

	// Inicia las dos goroutines que conforman el pipeline de streaming.
	go recibirFragmentos(stream, writer)
	go decodificarYReproducir(reader, ctx, done)
}

// recibirFragmentos se ejecuta en su propia goroutine y gestiona la capa de red.
// Su única responsabilidad es recibir fragmentos de audio del stream gRPC y pasarlos
// al pipe. Termina silenciosamente cuando el stream finaliza (io.EOF), se cancela,

// o si el 'PipeWriter' se cierra.
func recibirFragmentos(stream ss.AudioService_StreamAudioClient, writer *io.PipeWriter) {
	// Asegura que el writer del pipe se cierre al terminar la goroutine,
	// lo que señalará al reader que no llegarán más datos.
	defer writer.Close()
	for {
		fragmento, err := stream.Recv()
		if err != nil {
			// Cualquier error (fin de archivo, conexión cerrada, etc.) termina el bucle.
			return
		}

		if _, err := writer.Write(fragmento.GetData()); err != nil {
			// Si la escritura en el pipe falla (por ejemplo, porque el lector se cerró),
			// la goroutine también termina.
			return
		}
	}
}

// decodificarYReproducir se ejecuta en una goroutine y gestiona la capa de audio.
// Lee los datos del pipe, los decodifica y los reproduce. Utiliza un 'select' para
// manejar dos posibles eventos de finalización: que la canción termine naturalmente
// o que el usuario cancele la reproducción.
func decodificarYReproducir(reader *io.PipeReader, ctx context.Context, done chan bool) {
	// Este defer es crucial: garantiza que, sin importar cómo termine la función
	// (error, finalización, cancelación), siempre se cerrará el reader del pipe y
	// se notificará al canal 'done'.
	defer func() {
		reader.Close()
		done <- true
	}()

	// Intenta decodificar el stream de bytes entrante como un archivo MP3.
	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		// Este error es esperado si el contexto se cancela antes de recibir
		// suficientes datos para formar una cabecera MP3 válida.
		return
	}
	defer streamer.Close()

	// Inicializa el hardware de audio con el formato correcto la primera vez que se llama.
	speakerInicializado.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})

	// Crea un canal para recibir una señal cuando la reproducción de 'beep' termine.
	reproduccionTerminada := make(chan struct{})
	// Inicia la reproducción. beep.Seq ejecuta el callback cuando el 'streamer' termina.
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(reproduccionTerminada) // Cierra el canal para señalar el fin.
	})))

	// Espera concurrentemente por el primer evento que ocurra.
	select {
	case <-reproduccionTerminada:
		// La canción terminó de reproducirse de forma natural.
		fmt.Println("\nLa reproducción de audio ha finalizado.")
	case <-ctx.Done():
		// El contexto fue cancelado (el usuario detuvo la reproducción).
		// La función simplemente retorna; el 'defer' se encargará de la limpieza.
		return
	}
}
