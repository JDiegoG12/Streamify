package Controladores

import (
	// Esta ruta de importación ahora funcionará gracias al cambio en go.mod
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"fmt"
	"io"
	"log"
	"os"
)

type ServidorDeStreaming struct {
	ss.UnimplementedAudioServiceServer
}

func (s *ServidorDeStreaming) StreamAudio(req *ss.PeticionDTO, stream ss.AudioService_StreamAudioServer) error {
	nombreArchivo := req.GetTitulo()
	fmt.Printf("Petición remota recibida: StreamAudio para la canción: %s\n", nombreArchivo)

	rutaArchivo := "canciones/" + nombreArchivo + ".mp3"

	file, err := os.Open(rutaArchivo)
	if err != nil {
		log.Printf("Error al abrir el archivo %s: %v", rutaArchivo, err)
		return fmt.Errorf("no se pudo abrir el archivo de la canción: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 64*1024)
	fragmentoNum := 1

	for {
		bytesLeidos, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Canción enviada completamente desde el servidor.")
			break
		}
		if err != nil {
			log.Printf("Error leyendo el archivo: %v", err)
			return err
		}

		fmt.Printf("Fragmento #%d leído (%d bytes) y enviando...\n", fragmentoNum, bytesLeidos)

		if err := stream.Send(&ss.FragmentoCancion{Data: buffer[:bytesLeidos]}); err != nil {
			log.Printf("Error enviando fragmento al cliente: %v", err)
			return err
		}
		fragmentoNum++
	}

	return nil
}
