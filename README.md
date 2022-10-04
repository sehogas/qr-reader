# QR READER

Este proyecto lee QR desde cualquier cámara web, o cámara IP por RTSP. 

### Requerimientos previos

* Instalar https://chocolatey.org/install
* Instalar make con: choco install make ()
* Instalar cmake version 3.24.2
* Instalar mingw-w64: Buscar (https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)
* Instalar gocv, descargar fuentes y compilar.  
* zbar (depende)

### Instrucciones de instalación GOCV:

[Linux](https://gocv.io/getting-started/linux/)

[MacOs](https://gocv.io/getting-started/macos/)

[Windows](https://gocv.io/getting-started/windows/)

### Instrucciones de instalación ZBAR:

[Linux](http://zbar.sourceforge.net/download.html)

[MacOs](http://macappstore.org/zbar/)

[Windows](http://zbar.sourceforge.net/download.html)

```bash
go get github.com/bieber/barcode
```

Para comenzar, puedes clonar este reposito:
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