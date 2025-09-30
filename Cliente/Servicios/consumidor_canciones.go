package Servicios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"log"
)

func ListarGeneros(cliente sc.ServicioCancionesClient) []*sc.Genero {
	res, err := cliente.ListarGeneros(context.Background(), &sc.GetGenerosRequest{})
	if err != nil {
		log.Fatalf("Error al llamar a ListarGeneros: %v", err)
	}
	return res.GetGeneros()
}

func ListarCancionesPorGenero(cliente sc.ServicioCancionesClient, idGenero int32) []*sc.Cancion {
	req := &sc.GetCancionesPorGeneroRequest{IdGenero: idGenero}
	res, err := cliente.ListarCancionesPorGenero(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al llamar a ListarCancionesPorGenero: %v", err)
	}
	return res.GetCanciones()
}

func ConsultarCancion(cliente sc.ServicioCancionesClient, idCancion int32) *sc.Cancion {
	req := &sc.ConsultarCancionRequest{IdCancion: idCancion}
	res, err := cliente.ConsultarCancion(context.Background(), req)
	if err != nil {
		// En lugar de detener el programa con Fatalf, usamos Printf para que sea más robusto.
		log.Printf("Error al llamar a ConsultarCancion: %v", err)
		return nil
	}
	return res
}
