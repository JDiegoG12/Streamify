package main

import (
	"Streamify/Cliente/Vistas"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conexión con el Servidor de Canciones en el puerto 50052.
	connCanciones, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor de canciones: %v", err)
	}
	defer connCanciones.Close()
	clienteCanciones := sc.NewServicioCancionesClient(connCanciones)

	// Conexión con el Servidor de Streaming en el puerto 50051.
	connStreaming, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor de streaming: %v", err)
	}
	defer connStreaming.Close()
	clienteStreaming := ss.NewAudioServiceClient(connStreaming)

	// Inicia el menú principal de la aplicación.
	Vistas.MostrarMenuPrincipal(clienteCanciones, clienteStreaming)
}
