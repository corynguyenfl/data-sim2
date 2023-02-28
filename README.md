# Data Simulator

A simple program that publishes OpenFMB messages

## Configuration

The program is reading the following configuration file and generating messages according to the values.

- Set `pcc.is_closed` to `true` to close the breaker, and to `false` to open the breaker
- Set `w` to reflect the power measurement
- Set `is_on` to `true` to indicate that the device is in `StateKind_on`, and `false` to indicate that the device is in `StateKind_off`

```yaml
nats:
  url: 'nats://192.168.86.39:4222'
pcc:
  mrid: e6768784-48ad-40e9-af2a-9676413d4d6a
  w: 200.0
  is_closed: true  
ess:
  mrid: 836a8638-b448-4961-8258-47aa18e05f65 # control, status, reading
  reading_mrid: 6e595d68-67b4-434c-8c26-736104cc14fe # reading from Way 4
  soc: 65.0 
  is_on: true,
  mode: 2002 # valid mode: VSI_PQ: 2000, VSI_ISO: 2002   
solar:
  mrid: 540b292a-e600-4ae4-b077-40b892ae6970 
  reading_mrid: 081a9cd4-b7d2-4a7c-b089-42d0053f8e6a         
  is_on: True
  w: 774.93
generator:
  mrid: 8e202725-5a8d-45c2-b575-8161927c6770        
  reading_mrid: 8e202725-5a8d-45c2-b575-8161927c6770
  is_on: false
  w: 123.3
shop-meter:
  mrid: 0648ef71-cb63-4347-921a-9dbf178da687               
  w: 100.0      
load-bank:    
  mrid: 4cadfaed-4176-453d-9625-e9ae8ad70529 # control, status, reading
  reading_mrid: 2a09c597-d666-4354-a31a-dfe6da7727a9 # reading from ion meter        
  is_on: false
  w: 200.0
```

## Build

```
go build
```

### Run

```
APP_CONF=config/app.yaml go run .
```

## Docker

To build docker image:

```
docker-compose build
```

To run docker container:

```
docker-compose up
```