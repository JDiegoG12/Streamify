// El paquete main es el punto de entrada para la aplicación cliente de Streamify.
// Su responsabilidad principal es la inicialización y configuración de las conexiones
// con los servicios remotos antes de ceder el control a la capa de presentación (Vistas).
package main

import (
	"Streamify/Cliente/Vistas"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// main es la función principal y el punto de inicio de la ejecución del cliente.
// Su flujo de trabajo consiste en:
//  1. Establecer conexiones gRPC con los dos servidores backend: el Servidor de Canciones
//     y el Servidor de Streaming.
//  2. Crear los clientes gRPC (stubs) necesarios para interactuar con los servicios remotos.
//  3. Asegurar que las conexiones se cierren correctamente al finalizar el programa mediante 'defer'.
//  4. Iniciar la interfaz de usuario de la consola, pasando los clientes ya configurados.
func main() {
	// Establece una conexión gRPC con el Servidor de Canciones en el puerto 50052.
	connCanciones, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Si la conexión falla, el programa termina con un error fatal.
		log.Fatalf("No se pudo conectar al servidor de canciones: %v", err)
	}
	// Se programa el cierre de la conexión para cuando la función main termine.
	defer connCanciones.Close()
	clienteCanciones := sc.NewServicioCancionesClient(connCanciones)

	// Establece una conexión gRPC con el Servidor de Streaming en el puerto 50051.
	connStreaming, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Si la conexión falla, el programa termina con un error fatal.
		log.Fatalf("No se pudo conectar al servidor de streaming: %v", err)
	}
	// Se programa el cierre de la conexión para cuando la función main termine.
	defer connStreaming.Close()
	clienteStreaming := ss.NewAudioServiceClient(connStreaming)

	// Una vez que las conexiones y los clientes están listos, se inicia el menú principal
	// de la aplicación, que a partir de este punto gestionará toda la interacción con el usuario.
	Vistas.MostrarMenuPrincipal(clienteCanciones, clienteStreaming)
}
