package Vistas

import (
	"Streamify/Cliente/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/faiface/beep/speaker"
)

var reader = bufio.NewReader(os.Stdin)

// MostrarMenuPrincipal es el bucle principal de la aplicación.
func MostrarMenuPrincipal(clienteCanciones sc.ServicioCancionesClient, clienteStreaming ss.AudioServiceClient) {
	for {
		fmt.Println("\n===== Spotify =====")
		fmt.Println("1. Ver géneros")
		fmt.Println("2. Salir")
		fmt.Print("Seleccione una opción: ")

		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimSpace(opcion)

		switch opcion {
		case "1":
			menuVerGeneros(clienteCanciones, clienteStreaming)
		case "2":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// menuVerGeneros muestra la lista de géneros disponibles.
func menuVerGeneros(clienteCanciones sc.ServicioCancionesClient, clienteStreaming ss.AudioServiceClient) {
	generos := Fachada.ObtenerGeneros(clienteCanciones)
	for {
		fmt.Println("\n===== Géneros Disponibles =====")
		for i, genero := range generos {
			fmt.Printf("%d. %s\n", i+1, genero.Nombre)
		}
		fmt.Printf("%d. Atrás\n", len(generos)+1)
		fmt.Print("Seleccione un género: ")

		opcionStr, _ := reader.ReadString('\n')
		opcion, err := strconv.Atoi(strings.TrimSpace(opcionStr))

		if err != nil || opcion < 1 || opcion > len(generos)+1 {
			fmt.Println("Opción no válida.")
			continue
		}

		if opcion == len(generos)+1 {
			return
		}

		generoSeleccionado := generos[opcion-1]
		menuVerCanciones(clienteCanciones, clienteStreaming, generoSeleccionado)
	}
}

// menuVerCanciones muestra las canciones del género seleccionado.
func menuVerCanciones(clienteCanciones sc.ServicioCancionesClient, clienteStreaming ss.AudioServiceClient, genero *sc.Genero) {
	canciones := Fachada.ObtenerCanciones(clienteCanciones, genero.Id)
	for {
		fmt.Printf("\n===== Género: %s =====\n", genero.Nombre)
		for i, cancion := range canciones {
			fmt.Printf("%d. %s - %s\n", i+1, cancion.Artista, cancion.Titulo)
		}
		fmt.Printf("%d. Atrás\n", len(canciones)+1)
		fmt.Print("Seleccione una canción: ")

		opcionStr, _ := reader.ReadString('\n')
		opcion, err := strconv.Atoi(strings.TrimSpace(opcionStr))

		if err != nil || opcion < 1 || opcion > len(canciones)+1 {
			fmt.Println("Opción no válida.")
			continue
		}

		if opcion == len(canciones)+1 {
			return
		}
		cancionSeleccionada := canciones[opcion-1]
		menuDetalleCancion(clienteCanciones, clienteStreaming, cancionSeleccionada.Id)
	}
}

func menuDetalleCancion(clienteCanciones sc.ServicioCancionesClient, clienteStreaming ss.AudioServiceClient, idCancion int32) {
	// Realizamos la llamada remota para obtener los detalles frescos de la canción.
	fmt.Println("\nConsultando detalles de la canción...")
	cancion := Fachada.ObtenerDetalleCancion(clienteCanciones, idCancion)

	if cancion == nil {
		fmt.Println("No se pudieron obtener los detalles de la canción. Intente de nuevo.")
		return
	}

	fmt.Printf("\n===== Canción: %s - %s =====\n", cancion.Artista, cancion.Titulo)
	fmt.Printf(" Título de la canción: %s\n", cancion.Titulo)
	fmt.Printf(" Artista / Banda: %s\n", cancion.Artista)
	// El álbum no está en tu struct, si lo necesitas, debes agregarlo al .proto y al repositorio.
	fmt.Printf(" Año de lanzamiento: %d\n", cancion.AnioLanzamiento)
	fmt.Printf(" Duración: %s\n", cancion.Duracion)
	fmt.Println("\n1. Reproducir")
	fmt.Println("2. Atrás")
	fmt.Print("Seleccione una opción: ")

	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		reproducirConOpcionDeSalir(clienteStreaming, cancion)
	case "2":
		return
	default:
		fmt.Println("Opción no válida.")
	}
}

// Esta es la versión final y completa de la función de reproducción.
func reproducirConOpcionDeSalir(clienteStreaming ss.AudioServiceClient, cancion *sc.Cancion) {
	// 1. Crear un contexto que podamos cancelar.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Canal para saber cuándo las goroutines de fondo han terminado.
	done := make(chan bool)

	// 2. Iniciar el streaming en segundo plano.
	go Fachada.IniciarStreaming(clienteStreaming, cancion.Titulo, ctx, done)

	// 3. Mostrar el menú de reproducción.
	fmt.Printf("\n===== Spotify =====\n")
	fmt.Printf("Canción: %s - %s\n\n", cancion.Artista, cancion.Titulo)
	fmt.Println("  Reproduciendo canción...")
	fmt.Println("\n1. Salir")
	fmt.Print("Seleccione una opción: ")

	// 4. Iniciar una goroutine para leer la entrada del usuario sin bloquear.
	userInput := make(chan string)
	go func() {
		input, _ := reader.ReadString('\n')
		userInput <- strings.TrimSpace(input)
	}()

	// 5. Esperar a que la canción termine o a que el usuario quiera salir.
	select {
	case <-done:
		// La canción terminó por sí sola.
		return
	case input := <-userInput:
		if input == "1" {
			fmt.Println("\nDeteniendo reproducción...")
			// Secuencia de parada:
			// 1. Vaciar el buffer de audio para un silencio inmediato.
			speaker.Clear()
			// 2. Cancelar el contexto para detener las goroutines.
			cancel()
			// 3. Esperar la confirmación de que las goroutines han terminado.
			<-done
		}
	}
}
