# Laboratorio de Sistemas Distribuidos: Sistema de Streaming de Audio con gRPC

**Universidad del Cauca**  
**Laboratorio de Sistemas Distribuidos**

### Autores
*   Ana Sofia Arango Yanza
*   Juan Diego Gomez Garces

---

## Descripción del Proyecto

Este proyecto consiste en la implementación de un sistema distribuido para la consulta y reproducción de música, desarrollado en **Go (Golang)** y utilizando el framework **gRPC** para la comunicación entre componentes.

El objetivo principal es demostrar el funcionamiento del modelo de Llamada a Procedimiento Remoto (RPC) para construir un sistema cliente-servidor robusto. La aplicación permite a un usuario navegar por un catálogo de canciones y géneros musicales y reproducir una canción seleccionada mediante **streaming de audio en tiempo real**, sin necesidad de descargar el archivo por completo.

## Arquitectura del Sistema

El sistema está compuesto por tres componentes principales que se ejecutan de forma independiente:

1.  **Servidor de Canciones (`ServidorCanciones`):**
    *   **Responsabilidad:** Gestiona toda la información o *metadata* de las canciones y los géneros musicales (ID, título, artista, año, etc.).
    *   **Comunicación:** Expone servicios gRPC para que los clientes puedan solicitar la lista de géneros disponibles y la lista de canciones pertenecientes a un género específico.

2.  **Servidor de Streaming (`ServidorStreaming`):**
    *   **Responsabilidad:** Almacena los archivos de audio (en formato MP3) y se encarga de transmitirlos.
    *   **Comunicación:** Expone un servicio gRPC de tipo *server-streaming*. Cuando un cliente solicita una canción, este servidor la lee y la envía en pequeños fragmentos (chunks) consecutivos, permitiendo la reproducción inmediata en el cliente.

3.  **Cliente (`Cliente`):**
    *   **Responsabilidad:** Es una aplicación de consola que provee la interfaz de usuario. Permite al usuario interactuar con el sistema a través de menús.
    *   **Comunicación:** Se conecta a ambos servidores. Utiliza llamadas RPC unarias para obtener la metadata del Servidor de Canciones y establece una conexión de streaming con el Servidor de Streaming para recibir y reproducir el audio de una canción.

## Tecnologías y Librerías Clave

*   **Lenguaje de Programación:** Go (Golang)
*   **Framework de Comunicación:** gRPC
*   **Definición de Interfaces:** Protocol Buffers (Protobuf v3)
*   **Concurrencia:** Goroutines y Canales de Go para manejar el streaming de audio de forma asíncrona.
*   **Reproducción de Audio (Cliente):** Librería `github.com/faiface/beep` para la decodificación y reproducción de audio MP3.

## Patrones de Diseño Utilizados

Siguiendo los requerimientos, el proyecto se estructuró aplicando los siguientes patrones:

*   **Patrón en Capas:** La lógica en los servidores y el cliente está organizada para separar responsabilidades:
    *   **Controladores/Vistas:** Manejan las peticiones gRPC y la interacción con el usuario.
    *   **Fachada:** Simplifica el acceso a la lógica de negocio.
    *   **Acceso a Datos / Servicios:** Se encargan de la lógica de obtención de datos (en este caso, de repositorios en memoria) y la comunicación.
*   **DTO (Data Transfer Object):** Los mensajes definidos en los archivos `.proto` (como `Cancion`, `Genero`) actúan como DTOs, encapsulando de forma estandarizada los datos que viajan a través de la red.
*   **MVC (Modelo-Vista-Controlador):** La estructura del cliente sigue un enfoque similar a MVC, donde las `Vistas` gestionan la consola, la `Fachada` actúa como controlador y los `Modelos` son los DTOs que contienen los datos.

## ¿Cómo Compilar y Ejecutar el Proyecto?

### Prerrequisitos

1.  Tener instalado **Go** (versión 1.18 o superior).
2.  Tener instalado el compilador de Protocol Buffers (`protoc`).
3.  Tener instalados los plugins de `protoc` para Go y gRPC:
    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
    Asegúrate de que tu variable de entorno `GOPATH` esté configurada correctamente y que `$GOPATH/bin` esté en tu `PATH`.

### Pasos para la Ejecución

Es necesario abrir **3 terminales separadas**, una para cada componente del sistema.

**1. Terminal 1 - Iniciar Servidor de Canciones:**

```sh
# Navegar al directorio del servidor
cd ServidorCanciones

# Ejecutar el servidor
go run Main/servidor_canciones.go
