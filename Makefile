########################################################
override TARGET=qr-reader
VERSION=1.0
OS=windows
ARCH=amd64
FLAGS="-s -w"
CGO=1
CGO_LDFLAGS=-static -s
GOEXE=".exe"
NETWORK=sigep-network
########################################################

run:
	@echo Ejecutando programa...
	SET CC=x86_64-w64-mingw32-gcc
	SET AR=x86_64-w64-mingw32-ar
	SET CGO_ENABLED=$(CGO) 
	SET CGO_LDFLAGS=$(CGO_LDFLAGS)
	SET CGO_CXXFLAGS=-static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic
	SET GOEXE=$(GOEXE)
	SET GOOS=$(OS) 
	SET GOARCH=$(ARCH) 
	go run main.go --client-id "Prueba" --zone-id "AQ" --event-id "I" --db-name "data.db" --device-id 0

bin:
	@echo Generando binario ... (en windows con poweshell)
	SET CC=x86_64-w64-mingw32-gcc
	SET AR=x86_64-w64-mingw32-ar
	SET CGO_ENABLED=$(CGO) 
	SET CGO_LDFLAGS=$(CGO_LDFLAGS)
	SET CGO_CXXFLAGS=-static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic
	SET GOEXE=$(GOEXE)
	SET GOOS=$(OS) 
	SET GOARCH=$(ARCH) 
	go build -v -x -ldflags=$(FLAGS)  .

exec:
	./qr-reader --client-id "Prueba" --zone-id "AQ" --event-id "I" --db-name "data.db" --device-id 0

install: 
	@echo Instalando binario ... (en windows con poweshell)
	@echo CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH)  go install -ldflags=$(FLAGS) 
	@go install -tags sqlite_userauth -ldflags=$(FLAGS)
	SET CGO_ENABLED=$(CGO) 
	SET GOOS=$(OS) 
	SET GOARCH=$(ARCH) 
	go install -ldflags=$(FLAGS)
	@go install -ldflags=$(FLAGS)

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

cleanW:
	@echo Borrando binario ...
	del $(TARGET).exe
	del data.db

.PHONY: clean run install build start stop createnetwork start_with_network
.DEFAULT: 
	@echo 'No hay disponible ninguna regla para este destino'
