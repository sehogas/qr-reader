########################################################
override TARGET=qr-reader
VERSION=1.0
CC=x86_64-w64-mingw32-gcc
AR=x86_64-w64-mingw32-gcc-ar
CGO_ENABLED=1
CGO_LDFLAGS='-static -s'
CGO_CXXFLAGS='-static-libgcc -static-libstdc++ -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic'
GOEXE=.exe
GOOS=windows
GOARCH=amd64
NETWORK=sigep-network
FLAGS='-s -w'
########################################################

run: cam1

cam1:
	go run main.go run --file cam1_aq_i_x1.env.encrypted

cam2:
	go run main.go run --file cam2_aq_s_x1.env.encrypted

cam3:
	go run main.go run --file cam3_ap_i_x1.env.encrypted

cam4:
	go run main.go run --file cam4_aq_i_x2.env.encrypted

cam5:
	go run main.go run --file hikvision.env.encrypted

cam6:
	go run main.go run --file cam6_aq_i_x2.env.encrypted

cam7:
	go run main.go run --file cam7_os_i_x1.env.encrypted


encrypt1:
	go run main.go encrypt --file cam1_aq_i_x1.env

encrypt2:
	go run main.go encrypt --file cam2_aq_s_x1.env

encrypt3: 
	go run main.go encrypt --file cam3_ap_i_x1.env

encrypt4: 
	go run main.go encrypt --file cam4_aq_i_x2.env

encrypt5: 
	go run main.go encrypt --file hikvision.env

encrypt6: 
	go run main.go encrypt --file cam6_aq_i_x2.env

encrypt7: 
	go run main.go encrypt --file cam7_os_i_x1.env

bin:
	@echo Generando binario... 
	go build -v -x -ldflags=$(FLAGS)  .

exec1:
	@echo Ejecutando modo 1 QR...
	./qr-reader run --file=cam1_aq_i_x1.env.encrypted

exec2:
	@echo Ejecutando modo 2 QR...
	./qr-reader run --file=cam2_aq_s_x1.env.encrypted

exec3:
	@echo Ejecutando modo 1 QR...
	./qr-reader run --file=cam3_ap_i_x1.env.encrypted

exec4:
	@echo Ejecutando modo 2 QR...
	./qr-reader run --file=cam4_aq_i_x2.env.encrypted

exec7:
	@echo Ejecutando OFICINA SEGURIDAD modo 1 QR...
	./qr-reader run --file=cam7_os_i_x1.env.encrypted

install: 
	@echo Instalando binario...
	@echo CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH)  go install -ldflags=$(FLAGS) 
	@go install -tags sqlite_userauth -ldflags=$(FLAGS)
	SET CGO_ENABLED=$(CGO) 
	SET GOOS=$(OS) 
	SET GOARCH=$(ARCH) 
	go install -ldflags=$(FLAGS)
	@go install -ldflags=$(FLAGS)

build:
	@echo Construyendo imagen docker...
	docker build -t $(TARGET):$(VERSION) .
	docker tag $(TARGET):$(VERSION) $(TARGET):latest

start:
	@echo Ejecutando contenedor docker...
	docker run --rm -d --name $(TARGET) -p 3000:3000 $(TARGET):latest

start_with_network:
	@echo Ejecutando contenedor docker con network...
	docker run --rm -d --name $(TARGET) --network $(NETWORK) -p 3000:3000 $(TARGET):latest

stop:
	@echo Parando contenedor docker...
	docker stop $(TARGET)

createnetwork:
	@echo Creando red docker...
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
