# Dummy client unix socket

Cliente dummy para testear conexiones a sockets unix.

### Features
- Cantidad variable de conexiones.
- Tiempo entre lanzamientos entre conexiones, tiene la finalidad de no saturar la cola de conexiones en espera 
  del socket, sobre todo en kernels viejos donde “somaxconn” esta seteado por defecto en 128.
- No cerrar el socket luego de enviar los datos, tiene la finalidad de testear el servidor 
  con gran cantidad de files descriptors monitoreados al mismo tiempo
- Tiempo de espera entre que envió los datos y cierra la aplicación. 
  Esto da tiempo de revisar el estado delos sockets con la herramienta `lsof`

### Usage

```
Usage of ./sock_client:
    -f string
          Unix socket path (default "./echo.sock")
    -n uint
          Number of conexion (And goroutines) (default 200)
    -t uint
          Time before launch a new goroutines (in ms) (default 2)
    -u    Don't close socket after send data
    -w uint
          Wait time between the data was sent and the application closes (seconds)
```

