package main

import (
	"Streamify/ServidorCanciones/Controladores"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Iniciando Servidor de Canciones...")

	// El servidor escuchará en el puerto 50052.
	// El servidor de streaming usará el 50051.
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Fallo al escuchar en el puerto 50052: %v", err)
	}

	// Crea una nueva instancia del servidor gRPC.
	s := grpc.NewServer()

	// Registra nuestro controlador (que implementa la interfaz del servicio) con el servidor gRPC.
	sc.RegisterServicioCancionesServer(s, &Controladores.ServidorDeCanciones{})

	fmt.Println("Servidor de Canciones gRPC escuchando en el puerto :50052")

	// Inicia el servidor. Se quedará bloqueado aquí, esperando conexiones.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC: %v", err)
	}
}
