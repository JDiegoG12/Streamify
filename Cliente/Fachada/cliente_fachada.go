package Fachada

import (
	"Streamify/Cliente/Servicios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"context"
)

// ObtenerGeneros es una simple envoltura para la llamada al servicio.
func ObtenerGeneros(clienteCanciones sc.ServicioCancionesClient) []*sc.Genero {
	return Servicios.ListarGeneros(clienteCanciones)
}

// ObtenerCanciones es una envoltura para la llamada al servicio.
func ObtenerCanciones(clienteCanciones sc.ServicioCancionesClient, idGenero int32) []*sc.Cancion {
	return Servicios.ListarCancionesPorGenero(clienteCanciones, idGenero)
}

// Añadimos el parámetro de contexto
func IniciarStreaming(clienteStreaming ss.AudioServiceClient, titulo string, ctx context.Context, done chan bool) {
	// Y lo pasamos al servicio
	Servicios.ReproducirCancion(clienteStreaming, titulo, ctx, done)
}
