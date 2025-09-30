// El paquete Servicios encapsula la lógica de comunicación del cliente con los
// microservicios remotos. Este archivo, en particular, implementa los 'consumidores'
// que interactúan con los puntos finales (endpoints) del Servidor de Canciones,
// manejando la creación de peticiones gRPC y el procesamiento de sus respuestas.
package Servicios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"log"
)

// ListarGeneros realiza una llamada RPC para obtener la lista completa de géneros musicales
// disponibles en el servidor de canciones.
//
// Parámetros:
//   - cliente: El cliente gRPC (stub) ya inicializado para el servicio de canciones.
//
// Retorna:
//   - Un slice de punteros a 'sc.Genero', representando todos los géneros disponibles.
//
// Nota: En caso de un error en la comunicación, el programa terminará abruptamente,
// ya que se considera un fallo crítico para la funcionalidad de la aplicación.
func ListarGeneros(cliente sc.ServicioCancionesClient) []*sc.Genero {
	// Se utiliza context.Background() ya que esta es una operación de corta duración
	// y sin necesidad de cancelación explícita.
	res, err := cliente.ListarGeneros(context.Background(), &sc.GetGenerosRequest{})
	if err != nil {
		log.Fatalf("Error al llamar a ListarGeneros: %v", err)
	}
	return res.GetGeneros()
}

// ListarCancionesPorGenero realiza una llamada RPC para obtener las canciones asociadas
// a un género específico, identificado por su ID.
//
// Parámetros:
//   - cliente: El cliente gRPC (stub) del servicio de canciones.
//   - idGenero: El identificador único del género para filtrar la búsqueda.
//
// Retorna:
//   - Un slice de punteros a 'sc.Cancion' que corresponden al género solicitado.
//
// Nota: Al igual que ListarGeneros, un fallo en esta llamada se considera crítico
// y detendrá la ejecución del cliente.
func ListarCancionesPorGenero(cliente sc.ServicioCancionesClient, idGenero int32) []*sc.Cancion {
	req := &sc.GetCancionesPorGeneroRequest{IdGenero: idGenero}
	res, err := cliente.ListarCancionesPorGenero(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al llamar a ListarCancionesPorGenero: %v", err)
	}
	return res.GetCanciones()
}

// ConsultarCancion realiza una llamada RPC para obtener los detalles completos de una
// única canción a partir de su ID.
//
// Parámetros:
//   - cliente: El cliente gRPC (stub) del servicio de canciones.
//   - idCancion: El identificador único de la canción a consultar.
//
// Retorna:
//   - Un puntero a 'sc.Cancion' con la información detallada de la canción.
//   - Retorna 'nil' si la canción no se encuentra o si ocurre un error en la comunicación.
//
// Nota sobre el manejo de errores: A diferencia de las funciones de listado, esta
// función no termina el programa en caso de error. En su lugar, registra el error y
// devuelve 'nil', permitiendo que la capa de presentación (Vistas) maneje el fallo
// de forma controlada (por ejemplo, mostrando un mensaje al usuario sin cerrar la app).
func ConsultarCancion(cliente sc.ServicioCancionesClient, idCancion int32) *sc.Cancion {
	req := &sc.ConsultarCancionRequest{IdCancion: idCancion}
	// La llamada gRPC retorna el objeto de respuesta directamente, que coincide
	// con el tipo de dato de la canción.
	res, err := cliente.ConsultarCancion(context.Background(), req)
	if err != nil {
		log.Printf("Error al llamar a ConsultarCancion: %v", err)
		return nil
	}
	return res
}
