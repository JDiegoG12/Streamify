// El paquete Vistas implementa la capa de presentación (la "V" en un modelo MVC)
// de la aplicación cliente. Su responsabilidad es renderizar la interfaz de usuario
// en la consola, capturar las entradas del usuario y coordinar la navegación entre
// los diferentes menús. Utiliza el paquete Fachada para solicitar datos y ejecutar
// acciones, sin tener conocimiento directo de la comunicación gRPC.
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

// reader es un lector de bufer para la entrada estándar (os.Stdin), compartido
// a lo largo del paquete para leer las opciones del usuario de manera eficiente.
var reader = bufio.NewReader(os.Stdin)

// MostrarMenuPrincipal es el punto de entrada y el bucle principal de la interfaz de usuario.
// Presenta las opciones de nivel superior y dirige al usuario a las subsecciones
// correspondientes. El bucle se ejecuta indefinidamente hasta que el usuario elige salir.
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

// menuVerGeneros se encarga de mostrar la lista de géneros musicales disponibles.
// Llama a la fachada para obtener los datos, los presenta al usuario y gestiona la
// selección para navegar al menú de canciones o para volver al menú principal.
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

		// Validación de la entrada del usuario.
		if err != nil || opcion < 1 || opcion > len(generos)+1 {
			fmt.Println("Opción no válida.")
			continue
		}

		if opcion == len(generos)+1 {
			return // Vuelve al menú principal.
		}

		generoSeleccionado := generos[opcion-1]
		menuVerCanciones(clienteCanciones, clienteStreaming, generoSeleccionado)
	}
}

// menuVerCanciones muestra las canciones que pertenecen a un género previamente seleccionado.
// Al igual que otros menús, obtiene los datos a través de la fachada y maneja la navegación.
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

		// Validación de la entrada del usuario.
		if err != nil || opcion < 1 || opcion > len(canciones)+1 {
			fmt.Println("Opción no válida.")
			continue
		}

		if opcion == len(canciones)+1 {
			return // Vuelve al menú de géneros.
		}
		cancionSeleccionada := canciones[opcion-1]
		menuDetalleCancion(clienteCanciones, clienteStreaming, cancionSeleccionada.Id)
	}
}

// menuDetalleCancion muestra la información detallada de una canción seleccionada
// y ofrece las opciones de reproducirla o volver atrás.
func menuDetalleCancion(clienteCanciones sc.ServicioCancionesClient, clienteStreaming ss.AudioServiceClient, idCancion int32) {
	fmt.Println("\nConsultando detalles de la canción...")
	cancion := Fachada.ObtenerDetalleCancion(clienteCanciones, idCancion)

	// Manejo robusto en caso de que la canción no se encuentre o haya un error de red.
	if cancion == nil {
		fmt.Println("No se pudieron obtener los detalles de la canción. Intente de nuevo.")
		return
	}

	fmt.Printf("\n===== Canción: %s - %s =====\n", cancion.Artista, cancion.Titulo)
	fmt.Printf(" Título de la canción: %s\n", cancion.Titulo)
	fmt.Printf(" Artista / Banda: %s\n", cancion.Artista)
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
		return // Vuelve al menú de canciones.
	default:
		fmt.Println("Opción no válida.")
	}
}

// reproducirConOpcionDeSalir gestiona la sesión de reproducción de una canción.
// Esta función es un excelente ejemplo de manejo de concurrencia en una UI de consola:
//
//  1. Inicia la reproducción de audio en segundo plano (en goroutines gestionadas por la fachada).
//  2. Lanza una goroutine dedicada exclusivamente a esperar la entrada del usuario para no bloquear
//     el hilo principal.
//  3. Utiliza un `select` para esperar concurrentemente a uno de dos eventos: que la canción
//     termine por sí sola, o que el usuario decida intervenir.
func reproducirConOpcionDeSalir(clienteStreaming ss.AudioServiceClient, cancion *sc.Cancion) {
	// Se crea un contexto cancelable para poder detener las operaciones de streaming.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan bool) // Canal para señalar la finalización del streaming.

	// Inicia el streaming en una nueva goroutine.
	go Fachada.IniciarStreaming(clienteStreaming, cancion.Titulo, ctx, done)

	fmt.Printf("\n===== Spotify =====\n")
	fmt.Printf("Canción: %s - %s\n\n", cancion.Artista, cancion.Titulo)
	fmt.Println("  Reproduciendo canción...")
	fmt.Println("\n1. Salir")
	fmt.Print("Seleccione una opción: ")

	// Goroutine para leer la entrada del usuario de forma no bloqueante.
	userInput := make(chan string)
	go func() {
		input, _ := reader.ReadString('\n')
		select {
		case userInput <- strings.TrimSpace(input):
		default: // El receptor ya no está esperando, la goroutine termina.
		}
	}()

	// Espera por el primer evento que ocurra.
	select {
	case <-done:
		// Caso 1: La canción termina de reproducirse naturalmente.
		// La goroutine que lee 'userInput' sigue bloqueada. Se solicita al usuario
		// presionar Enter para desbloquearla y asegurar una finalización limpia,
		// evitando una 'goroutine huérfana'.
		fmt.Println("\nPresione Enter para continuar...")
		cancel()
		<-userInput // Drena el canal para que la goroutine de entrada termine.
		return

	case input := <-userInput:
		// Caso 2: El usuario interrumpe la reproducción.
		if input == "1" {
			fmt.Println("\nDeteniendo reproducción...")
			// Secuencia de detención controlada:
			// 1. Limpiar el búfer de audio para un silencio inmediato.
			speaker.Clear()
			// 2. Cancelar el contexto para señalar a las goroutines de streaming que deben parar.
			cancel()
			// 3. Esperar la confirmación de que han terminado ('done').
			<-done
		}
		// Si el usuario ingresa algo distinto de "1", la función simplemente termina.
		// El 'defer cancel()' se ejecutará, deteniendo la música al salir.
	}
}
