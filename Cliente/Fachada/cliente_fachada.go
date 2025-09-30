// El paquete Fachada actúa como un intermediario simplificado entre la capa de
// presentación (Vistas) y la capa de servicios del cliente. Su propósito es
// desacoplar la lógica de la interfaz de usuario de los detalles de cómo se
// realizan las llamadas a los servicios remotos, siguiendo el patrón de diseño Facade.
package Fachada

import (
	"Streamify/Cliente/Servicios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"context"
)

// ObtenerGeneros solicita la lista completa de géneros musicales disponibles
// al servidor de canciones.
//
// Parámetros:
//   - clienteCanciones: El cliente gRPC para comunicarse con el servicio de canciones.
//
// Retorna:
//   - Un slice de punteros a objetos 'Genero', representando todos los géneros registrados.
func ObtenerGeneros(clienteCanciones sc.ServicioCancionesClient) []*sc.Genero {
	return Servicios.ListarGeneros(clienteCanciones)
}

// ObtenerCanciones solicita la lista de canciones que pertenecen a un género específico.
//
// Parámetros:
//   - clienteCanciones: El cliente gRPC para comunicarse con el servicio de canciones.
//   - idGenero: El identificador único del género cuyas canciones se desean obtener.
//
// Retorna:
//   - Un slice de punteros a objetos 'Cancion' que pertenecen al género especificado.
func ObtenerCanciones(clienteCanciones sc.ServicioCancionesClient, idGenero int32) []*sc.Cancion {
	return Servicios.ListarCancionesPorGenero(clienteCanciones, idGenero)
}

// ObtenerDetalleCancion consulta y retorna la información detallada de una canción específica
// a partir de su ID.
//
// Parámetros:
//   - clienteCanciones: El cliente gRPC para comunicarse con el servicio de canciones.
//   - idCancion: El identificador único de la canción que se desea consultar.
//
// Retorna:
//   - Un puntero al objeto 'Cancion' con todos sus detalles.
func ObtenerDetalleCancion(clienteCanciones sc.ServicioCancionesClient, idCancion int32) *sc.Cancion {
	return Servicios.ConsultarCancion(clienteCanciones, idCancion)
}

// IniciarStreaming comienza el proceso de reproducción de una canción. Esta función
// no bloquea la ejecución, sino que delega la operación de streaming a una goroutine
// en la capa de servicios.
//
// Parámetros:
//   - clienteStreaming: El cliente gRPC para comunicarse con el servicio de streaming.
//   - titulo: El título de la canción que se va a reproducir.
//   - ctx: El contexto que permite la cancelación de la operación de streaming (por ejemplo,
//     cuando el usuario decide detener la canción).
//   - done: Un canal booleano que recibirá una señal cuando el proceso de streaming
//     (ya sea por finalización o por cancelación) haya terminado por completo.
func IniciarStreaming(clienteStreaming ss.AudioServiceClient, titulo string, ctx context.Context, done chan bool) {
	Servicios.ReproducirCancion(clienteStreaming, titulo, ctx, done)
}
