# freetalkbot
Your own IA to handle communications with your customers via phonecalls or whatsapp

## Dependencies

* Go version recommended: 1.22
* [whatsapp-media-decrypt](https://github.com/ddz/whatsapp-media-decrypt/tree/master) tool
* Install dependencies with `go mod tidy`

For better golang developer experience you can install [golang-air](https://github.com/cosmtrek/air)

### Environment variables

Check the variables in `env.example` file. Create `.env` file with `cp -a .env.example .env` and modify it with your values. 
All bot servers reads configuration from this file. 

## Development

The components added in the `docker-compose.yml` file are:

* Asterisk
* [Rasa assistant](https://rasa.com/)

Raise up the development environment with `docker-compose up -d`

### Configure asterisk

1. Once raised up, copy local-config to container-config `cp -a asterisk/local-config/* asterisk/container-config/`
2. Restart asterisk container `docker-compose restart asterisk`

Asterisk is raised up in network_mode host. The asterisk configuration files are mapped in folder `asterisk/container-config`

### Register SIP endpoint

Checkout pjsip_endpoint.conf file in `asterisk/container-config` folder.

### Rasa assistant

The assistant is configured to be a reminderbot, inspired in the [example](https://github.com/RasaHQ/rasa/tree/main/examples/reminderbot) provided by rasa. The files in this folder are for NLP training of the assistant.
Checkout the following resources to get more knowledge about RASA.

* [Documentation](https://rasa.com/docs/rasa/training-data-format)
* [Youtube Channel](https://www.youtube.com/@RasaHQ)

If you modify any of the rasa files you will need to retrain the assistant, you can do it with `make rasa-train`

## Audio bot

```sh
go run main.go init -c audio
```

with air

```sh
make run-audio
```

## Whatsapp bot

```sh
go run main.go init -c whatsapp
```

with air

```sh
make run-whatsapp
```

After initialize you will see in the logs a QR code. Scan that QR code with the whatsapp account that you will use.
If you can't scan the QR code you can also link the whatsapp account using a pair code. For that you must set the envar `PAIR_PHONE_NUMBER` with 
your phone number using format show in the `.env.example`. If you don't need the pair code don't set this envar.
The whatsapp server store session in sqlite, so you will see a `.db` file. If you delete this file you will have to login using a new QR code.

### Features

Now the code is prepared to receive text or voice messages. If you want to make assistant based on text you should modify it.
The assistant on this repository is a `reminderbot` that will send you reminders based on voice messages that you are sending to it 