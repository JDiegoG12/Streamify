package Fachada

import (
	"Streamify/Cliente/Servicios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
)

// ObtenerGeneros es una simple envoltura para la llamada al servicio.
func ObtenerGeneros(clienteCanciones sc.ServicioCancionesClient) []*sc.Genero {
	return Servicios.ListarGeneros(clienteCanciones)
}

// ObtenerCanciones es una envoltura para la llamada al servicio.
func ObtenerCanciones(clienteCanciones sc.ServicioCancionesClient, idGenero int32) []*sc.Cancion {
	return Servicios.ListarCancionesPorGenero(clienteCanciones, idGenero)
}

// IniciarStreaming es una envoltura para la llamada al servicio de streaming.
func IniciarStreaming(clienteStreaming ss.AudioServiceClient, titulo string) {
	Servicios.ReproducirCancion(clienteStreaming, titulo)
}
