package Vistas

import (
	"Streamify/Cliente/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		menuDetalleCancion(clienteStreaming, cancionSeleccionada)
	}
}

// menuDetalleCancion muestra los detalles de una canción y permite reproducirla.
func menuDetalleCancion(clienteStreaming ss.AudioServiceClient, cancion *sc.Cancion) {
	for {
		fmt.Printf("\n===== Canción: %s - %s =====\n", cancion.Artista, cancion.Titulo)
		fmt.Printf(" Título de la canción: %s\n", cancion.Titulo)
		fmt.Printf(" Artista / Banda: %s\n", cancion.Artista)
		fmt.Printf(" Año de lanzamiento: %d\n", cancion.AnioLanzamiento)
		fmt.Printf(" Duración: %s\n", cancion.Duracion)
		fmt.Println("\n1. Reproducir")
		fmt.Println("2. Salir")
		fmt.Print("Seleccione una opción: ")

		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimSpace(opcion)

		switch opcion {
		case "1":
			Fachada.IniciarStreaming(clienteStreaming, cancion.Titulo)
			// Espera a que el usuario presione Enter para continuar
			reader.ReadString('\n')
		case "2":
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}
