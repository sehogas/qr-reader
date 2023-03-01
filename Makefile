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

run: run_cam1

run_cam1_prod:
	go run main.go run --file prod.env.encrypted

run_cam1:
	go run main.go run --file cam1_aq_i.env.encrypted

run_cam2:
	go run main.go run --file cam2_aq_s.env.encrypted

run_cam3:
	go run main.go run --file cam3_ap_i.env.encrypted

run_cam4:
	go run main.go run --file cam4_aq_i_x2.env.encrypted

run_cam5:
	go run main.go run --file hikvision.env.encrypted

run_cam6:
	go run main.go run --file cam6_aq_i_x2.env.encrypted

encrypt1:
	@echo Ejecutando programa...
	go run main.go encrypt --file cam1_aq_i.env

encrypt2:
	@echo Ejecutando programa...
	go run main.go encrypt --file cam2_aq_s.env

runencrypt3: 
	go run main.go encrypt --file cam3_ap_i.env

encrypt3: sets runencrypt3

runencrypt4: 
	go run main.go encrypt --file cam4_aq_i_x2.env

encrypt4: sets runencrypt4

runencrypt5: 
	go run main.go encrypt --file hikvision.env

encrypt5: sets runencrypt5

runencrypt6: 
	go run main.go encrypt --file cam6_aq_i_x2.env


encrypt6: sets runencrypt6

runbin:
	@echo Generando binario ... en windows con poweshell
	go build -v -x -ldflags=$(FLAGS)  .

bin: sets runbin

exec1:
	@echo Ejecutando modo 1 QR ... [en windows con poweshell]
	./qr-reader run --file=cam1_aq_i.env.encrypted

exec2:
	@echo Ejecutando modo 2 QR ... [en windows con poweshell]
	./qr-reader run --file=cam2_aq_s.env.encrypted

exec3:
	@echo Ejecutando modo 1 QR ... [en windows con poweshell]
	./qr-reader run --file=cam3_ap_i.env.encrypted

install: 
	@echo Instalando binario ... [en windows con poweshell]
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
