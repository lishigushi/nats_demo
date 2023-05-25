## run contaners
``` 
docker compose up -d
```

## stop containers
```
docker compose down
```

## for leaf setup

run all the containers, then use the script for leaf stream setup

for the last question just press the Enter key
```
./setup_leaf.sh

? Adjust mirror start (y/N) n
? Import mirror from a different JetStream domain (y/N) y
? Foreign JetStream domain name [? for help] main
? Delivery prefix [? for help]
```

## data

All temporary data will be stored in the ./data subfolder.

## grafana

http://localhost:3000

for grafana access use default admin:admin

## publisher

first install dependencies
```
cd ./publisher
pip install -r requirements.txt
```

then just run (by default it publishes 100000 messages)
```
python3 publisher.py
```

If you need to specify the number of messages, use the `--messages` parameter
```
python3 publisher.py --messages 1000
```

## subscriber

build
```
cd subscriber && make build
```

When launched without parameters, an ephemeral client will be used, and the stream will be obtained from the first message
```
./subscriber
```

If a durable client is needed, use the `--consumer` parameter
```
./subscriber --consumer myconsumername
```

if you need to start with a specific sequence ID, use the `--start-seq` parameter
```
./subscriber --start-seq 1000
```

## nats embedded

build
```
cd nats_embedded
make build
```

run
```
./embedded
```
