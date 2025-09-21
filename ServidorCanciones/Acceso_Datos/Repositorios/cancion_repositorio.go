package Repositorios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// Simulamos una tabla de canciones en una base de datos.
var canciones = []*sc.Cancion{
	// Rock (Genero ID: 3)
	{Id: 1, Titulo: "De musica ligera", Artista: "Soda Stereo", AnioLanzamiento: 1990, Duracion: "3:33", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 2, Titulo: "Entre dos tierras", Artista: "Heroes del Silencio", AnioLanzamiento: 1990, Duracion: "6:09", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 3, Titulo: "El baile de los que sobran", Artista: "Los Prisioneros", AnioLanzamiento: 1986, Duracion: "5:44", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 4, Titulo: "Lamento boliviano", Artista: "Enanitos Verdes", AnioLanzamiento: 1994, Duracion: "3:42", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 5, Titulo: "Rayando el sol", Artista: "Maná", AnioLanzamiento: 1990, Duracion: "4:10", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},

	// Salsa (Genero ID: 1)
	{Id: 6, Titulo: "Pedro Navaja", Artista: "Ruben Blades", AnioLanzamiento: 1978, Duracion: "7:22", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},
	{Id: 7, Titulo: "Lloraras", Artista: "Oscar D'León", AnioLanzamiento: 1975, Duracion: "3:40", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},

	// Cumbia (Genero ID: 2)
	{Id: 8, Titulo: "Como te voy a olvidar", Artista: "Los Ángeles Azules", AnioLanzamiento: 1996, Duracion: "4:28", Genero: &sc.Genero{Id: 2, Nombre: "Cumbia"}},
}

// ObtenerCancionesPorGenero busca en la lista de canciones y devuelve solo las que coinciden con el ID del género.
func ObtenerCancionesPorGenero(idGenero int32) []*sc.Cancion {
	var cancionesFiltradas []*sc.Cancion
	for _, cancion := range canciones {
		if cancion.Genero.Id == idGenero {
			cancionesFiltradas = append(cancionesFiltradas, cancion)
		}
	}
	return cancionesFiltradas
}
