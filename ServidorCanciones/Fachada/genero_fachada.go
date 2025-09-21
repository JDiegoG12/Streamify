package Fachada

import (
	"Streamify/ServidorCanciones/Acceso_Datos/Repositorios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// ObtenerTodosLosGeneros simplemente llama al repositorio para obtener los datos.
// En un proyecto más complejo, aquí podría haber lógica adicional.
func ObtenerTodosLosGeneros() []*sc.Genero {
	return Repositorios.ObtenerGeneros()
}
