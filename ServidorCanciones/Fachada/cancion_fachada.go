package Fachada

import (
	"Streamify/ServidorCanciones/Acceso_Datos/Repositorios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// ObtenerCancionesPorIdGenero llama al repositorio correspondiente.
func ObtenerCancionesPorIdGenero(idGenero int32) []*sc.Cancion {
	return Repositorios.ObtenerCancionesPorGenero(idGenero)
}
