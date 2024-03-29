# Data Simulator

A simple program that publishes OpenFMB messages

## Configuration

The program is reading the following configuration file and generating messages according to the values.

- Set `pcc.is_closed` to `true` to close the breaker, and to `false` to open the breaker
- Set `W` to reflect the power measurement
- Set `is_on` to `true` to indicate that the device is in `StateKind_on`, and `false` to indicate that the device is in `StateKind_off`

```yaml
log-message-enabled: false
nats:
    url: nats://192.168.86.39:4222
microgrid-controller:
    enabled: false
    pcc:
        mrid: e6768784-48ad-40e9-af2a-9676413d4d6a
        W: 0
        is_closed: true
    ess:
        mrid: 836a8638-b448-4961-8258-47aa18e05f65
        reading_mrid: 6e595d68-67b4-434c-8c26-736104cc14fe
        soc: 65
        mode: 2000
        is_on: true
        W: 0
    solar:
        mrid: 540b292a-e600-4ae4-b077-40b892ae6970
        is_on: true
        w: 0
    generator:
        mrid: 8e202725-5a8d-45c2-b575-8161927c6770
        reading_mrid: 8e202725-5a8d-45c2-b575-8161927c6770
        is_on: false
        W: 0
    shop-meter:
        mrid: 0648ef71-cb63-4347-921a-9dbf178da687
        W: 0
    load-bank:
        mrid: 4cadfaed-4176-453d-9625-e9ae8ad70529
        reading_mrid: 2a09c597-d666-4354-a31a-dfe6da7727a9
        is_on: false
        W: 0
cvr:
    enabled: true
    recloser1:
        name: recloser1
        mrid: 90647973-e0cc-40c8-ad3c-b8adf6c7db05
        Va: 122
        Vb: 103
        Vc: 100
        W: 0
        is_closed: false
    recloser2:
        name: recloser2
        mrid: 9fd20ac0-35b7-4b8f-85e1-eaa8e771c6fd
        Va: 102
        Vb: 103
        Vc: 100
        W: 0
        is_closed: false
    vr1:
        name: vr1
        mrid: 82ad21e3-6f86-437b-a360-405ea3dca012
        pos: -5
        volLmHi: false
        volLmLo: false
        voltageSetPointEnabled: true
        source_primary_voltage: 123
        source_secondary_voltage: 123
        load_primary_voltage: 123
        load_secondary_voltage: 123
    vr2:
        name: vr2
        mrid: 953ef053-c288-463a-a2ce-33c021ae61e9
        pos: -5
        volLmHi: false
        volLmLo: false
        voltageSetPointEnabled: true
        source_primary_voltage: 123
        source_secondary_voltage: 123
        load_primary_voltage: 123
        load_secondary_voltage: 123
    vr3:
        name: vr3
        mrid: 9be8b771-7cd5-4335-95f8-165d0aff04fc
        pos: 0
        volLmHi: false
        volLmLo: false
        voltageSetPointEnabled: true
        source_primary_voltage: 123
        source_secondary_voltage: 123
        load_primary_voltage: 123
        load_secondary_voltage: 123
    capbank:
        name: capbank
        mrid: ea6c74bd-2af3-450d-b459-b4bdf22371a3
        control-mode: 4
        is_closed: false
        volLmt: false
        varLmt: false
        tempLmt: false
        Ia: 0
        Ib: 0
        Ic: 0
        Va: 100
        Vb: 100
        Vc: 100
        V2a: 0
        V2b: 0
        V2c: 0
        Wa: 0
        Wb: 0
        Wc: 0
    load1:
        name: load1
        mrid: f4be354c-cdde-4050-84de-2ec3c19f6d70
        Ia: 10
        Ib: 11
        Ic: 12
        Va: 50
        Vb: 50
        Vc: 50
        Apparent: 100
        Reactive: 100
        W: 100
    load2:
        name: load2
        mrid: 5c6a381b-ea49-4c12-84f5-93d5defa3394
        Ia: 10
        Ib: 11
        Ic: 12
        Va: 50
        Vb: 50
        Vc: 50
        Apparent: 100
        Reactive: 100
        W: 100
    load3:
        name: load3
        mrid: 5c6a381b-ea49-4c12-84f5-93d5defa3394
        Ia: 10
        Ib: 11
        Ic: 12
        Va: 50
        Vb: 50
        Vc: 50
        Apparent: 100
        Reactive: 100
        W: 100
    load4:
        name: load4
        mrid: 5c6a381b-ea49-4c12-84f5-93d5defa3394
        Ia: 10
        Ib: 11
        Ic: 12
        Va: 50
        Vb: 50
        Vc: 50
        Apparent: 100
        Reactive: 100
        W: 100
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
