// El paquete main es el punto de entrada para el microservicio del Servidor de Streaming.
// Su única responsabilidad es realizar la configuración inicial del entorno, inicializar
// el servidor gRPC, registrar los servicios que proveerá y ponerlo a la escucha de
// peticiones de los clientes.
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

// La función 'main' arranca el microservicio del Servidor de Streaming. El proceso sigue
// los pasos estándar para un servidor gRPC, con una verificación de entorno adicional:
//  1. Verificar y preparar el sistema de archivos (la carpeta 'canciones').
//  2. Configurar el punto de escucha de red en un puerto TCP específico.
//  3. Crear la instancia del servidor gRPC.
//  4. Registrar la implementación del servicio (el controlador) con el servidor.
//  5. Iniciar el bucle de servicio para aceptar conexiones indefinidamente.
func main() {
	// Se añade una verificación de robustez. Para evitar que el servidor falle si el
	// directorio 'canciones' no existe, se comprueba su existencia y se crea si es
	// necesario. Esto mejora la experiencia del usuario al desplegar el servicio.
	if _, err := os.Stat("canciones"); os.IsNotExist(err) {
		fmt.Println("Creando directorio 'canciones'. Por favor, coloca tus archivos MP3 aquí.")
		os.Mkdir("canciones", 0755)
	}

	fmt.Println("Iniciando Servidor de Streaming...")

	// Se asigna el puerto 50051 a este servicio, distinto al del Servidor de Canciones,
	// siguiendo el principio de un puerto por microservicio.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		// Si el puerto está ocupado, la aplicación no puede funcionar y termina.
		log.Fatalf("Fallo al escuchar en el puerto 50051: %v", err)
	}

	// Crea una instancia 'vacía' del servidor gRPC, lista para ser configurada.
	s := grpc.NewServer()

	// Paso de registro: Se le indica al servidor 's' que cualquier llamada al 'AudioService'
	// debe ser dirigida a una instancia de nuestro 'Controladores.ServidorDeStreaming'.
	// Esto vincula la definición del servicio .proto con la implementación en Go.
	ss.RegisterAudioServiceServer(s, &Controladores.ServidorDeStreaming{})

	fmt.Println("Servidor de Streaming gRPC escuchando en el puerto :50051")

	// Inicia el servidor en un bucle infinito de escucha. Esta es una operación bloqueante;
	// el programa se detendrá en esta línea, esperando y manejando conexiones de clientes
	// hasta que el proceso sea terminado.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC de streaming: %v", err)
	}
}
