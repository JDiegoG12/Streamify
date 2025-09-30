// El paquete main es el punto de entrada para la aplicación del Servidor de Canciones.
// Su única responsabilidad es configurar e iniciar el servidor gRPC, registrar los
// servicios que ofrecerá y ponerlo a la escucha de peticiones entrantes.
package main

import (
	"Streamify/ServidorCanciones/Controladores"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// main es la función que arranca el microservicio del Servidor de Canciones.
// El proceso de arranque sigue los pasos estándar para un servidor gRPC:
//  1. Definir y escuchar en un puerto TCP específico.
//  2. Crear una nueva instancia del servidor gRPC.
//  3. Registrar las implementaciones de los servicios (en este caso, nuestro
//     'ServidorDeCanciones' del paquete Controladores) con la instancia del servidor.
//  4. Iniciar el servidor para que comience a aceptar y procesar peticiones de los clientes.
func main() {
	fmt.Println("Iniciando Servidor de Canciones...")

	// Se especifica el puerto 50052 para este servicio. Es una práctica común en
	// arquitecturas de microservicios asignar puertos fijos y conocidos a cada servicio.
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		// Si el puerto ya está en uso o hay otro problema de red, el programa
		// no puede continuar y termina con un error fatal.
		log.Fatalf("Fallo al escuchar en el puerto 50052: %v", err)
	}

	// Se crea la instancia del servidor gRPC. En este punto, el servidor aún no
	// sabe qué servicios ofrecer ni cómo manejar las peticiones.
	s := grpc.NewServer()

	// Este es el paso clave de registro. Aquí se le dice al servidor 's' que todas
	// las peticiones que lleguen para el 'ServicioCanciones' deben ser manejadas
	// por una instancia de nuestro 'Controladores.ServidorDeCanciones'.
	// Esto conecta efectivamente nuestras implementaciones de los métodos RPC con
	// el motor de gRPC.
	sc.RegisterServicioCancionesServer(s, &Controladores.ServidorDeCanciones{})

	fmt.Println("Servidor de Canciones gRPC escuchando en el puerto :50052")

	// 's.Serve(lis)' es una operación bloqueante. Pone en marcha el servidor, que
	// comienza a escuchar peticiones en el 'listener' TCP. El programa se detendrá
	// en esta línea, procesando peticiones indefinidamente hasta que el proceso
	// sea terminado externamente.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC: %v", err)
	}
}
