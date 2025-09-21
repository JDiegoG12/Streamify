package main

import (
	"Streamify/ServidorStreaming/Controladores"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	// Verificamos si la carpeta 'canciones' existe. Si no, la creamos.
	if _, err := os.Stat("canciones"); os.IsNotExist(err) {
		fmt.Println("Creando directorio 'canciones'. Por favor, coloca tus archivos MP3 aquí.")
		os.Mkdir("canciones", 0755)
	}

	fmt.Println("Iniciando Servidor de Streaming...")

	// Usaremos el puerto 50051 para este servidor.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar en el puerto 50051: %v", err)
	}

	// Creamos la instancia del servidor gRPC.
	s := grpc.NewServer()

	// Registramos nuestro controlador de streaming.
	ss.RegisterAudioServiceServer(s, &Controladores.ServidorDeStreaming{})

	fmt.Println("Servidor de Streaming gRPC escuchando en el puerto :50051")

	// Iniciamos el servidor.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC de streaming: %v", err)
	}
}
