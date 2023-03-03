# QR READER

Este proyecto es un sistema autómata por línea de comando que lee QRs desde cualquier cámara web, o cámara IP por RTSP, validando que exista y sea valido en su base de datos local. Además, sincroniza cada cierto tiempo sus que información local contra un backend. 

### Requerimientos previos

* Instalar https://chocolatey.org/install
* Instalar make con: choco install make ()
* Instalar cmake version 3.24.2
* Instalar mingw-w64: Buscar (https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)
* Instalar gocv, descargar fuentes y compilar
* Requiere un 

### Instrucciones de instalación GOCV:

[Linux](https://gocv.io/getting-started/linux/)

[MacOs](https://gocv.io/getting-started/macos/)

[Windows](https://gocv.io/getting-started/windows/)


Para comenzar, puedes clonar este repositorio:
```bash
git clone https://github.com/sehogas/qr-reader.git qr-reader
cd qr-reader
go run main.go
```
O especificar dispositivo id (webcam):
```bash
go run main.go --device-id 1234
```

O especificar por ejemplo el flujo de datos (cámara ip):
```bash
go run main.go --from-file rtsp://user:pass@host:port/stream1
```

### Nota:
Ejecutar el programa sin parámetros intentará abrir su webcam predeterminada correspondiente al device id: 0