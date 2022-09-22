########################################################
override TARGET=qr-reader
VERSION=1.0
OS=linux
ARCH=windows
FLAGS="-s -w"
CGO=0
NETWORK=sigep-network
########################################################

run:
	@echo Ejecutando programa...
	go run main.go

bin:
	@echo Generando binario ...
	CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags=$(FLAGS) -o $(TARGET) .

install: 
	@echo Instalando binario ...
	SET CGO_ENABLED=$(CGO) 
	SET GOOS=$(OS) 
	SET GOARCH=$(ARCH) 
	go install -ldflags=$(FLAGS) 

build:
	@echo Construyendo imagen docker $(TARGET):$(VERSION) ...
	docker build -t $(TARGET):$(VERSION) .
	docker tag $(TARGET):$(VERSION) $(TARGET):latest

start:
	@echo Ejecutando contenedor docker $(TARGET):$(VERSION) ...
	docker run --rm -d --name $(TARGET) -p 3000:3000 $(TARGET):latest

start_with_network:
	@echo Ejecutando contenedor docker $(TARGET):$(VERSION) ...
	docker run --rm -d --name $(TARGET) --network $(NETWORK) -p 3000:3000 $(TARGET):latest

stop:
	docker stop $(TARGET)

createnetwork:
	docker network create -d bridge $(NETWORK)

#swagger:
#	swag init

#desa:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backdesa_sigep .
#	scp backdesa_sigep shogas@192.168.3.7:/home/shogas/go/bin

clean:
	@echo Borrando binario ...
	rm -rf $(TARGET)

.PHONY: clean run install build start stop createnetwork start_with_network
.DEFAULT: 
	@echo 'No hay disponible ninguna regla para este destino'
