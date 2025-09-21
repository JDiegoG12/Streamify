package Repositorios

import (
	// Importamos el paquete generado por protoc para usar los structs (Genero, Cancion, etc.)
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// Simulamos una tabla de géneros en una base de datos.
var generos = []*sc.Genero{
	{Id: 1, Nombre: "Salsa"},
	{Id: 2, Nombre: "Cumbia"},
	{Id: 3, Nombre: "Rock"},
}

// ObtenerGeneros devuelve la lista completa de géneros.
func ObtenerGeneros() []*sc.Genero {
	return generos
}
