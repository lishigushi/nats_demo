import argparse
import asyncio
import logging
import uuid

import nats
from nats.js import JetStreamContext
from nats.js.errors import NotFoundError


async def main():
    args = parse_args()
    nats_url = make_nats_url(args.url)

    jc = await get_jetstream_context(nats_url)
    await create_stream_if_need(jc, args.subject, args.stream, 3)
    await publish_messages(jc, args.subject, args.messages)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--url",
                        type=str,
                        default="nats://localhost:4222,nats://localhost:4223,nats://localhost:4224",
                        help="nats url")
    parser.add_argument("--messages", type=int, default=100000, help="number of messages")
    parser.add_argument("--stream", type=str, default="teststream", help="nats stream name")
    parser.add_argument("--subject", type=str, default="testsubject", help="messages subject")
    args = parser.parse_args()
    return args


def make_nats_url(raw: str) -> list[str]:
    url = [x.strip() for x in raw.split(',')]
    return url


async def get_jetstream_context(nats_url: list[str], domain=None) -> JetStreamContext:
    nc = await nats.connect(nats_url)
    jc = nc.jetstream(domain=domain) if domain else nc.jetstream()
    return jc


async def create_stream_if_need(jc: JetStreamContext, subject: str, stream: str, num_replicas: int) -> None:
    try:
        await jc.stream_info(stream)
        logging.info(f"stream {stream} exists already")
    except NotFoundError:
        try:
            await jc.add_stream(
                name=stream,
                subjects=[subject],
                num_replicas=num_replicas,
            )
            logging.info(f"stream {stream} created successfully")
        except Exception as e:
            logging.error(f"Failed to create stream {stream}. Error: {e}")


async def publish_messages(jc: JetStreamContext, subject: str, num_messages: int) -> None:
    for i in range(num_messages):
        await jc.publish(subject, str(uuid.uuid4()).encode())
        if i % 1000 == 0:
            print(f"published {i} messages")


if __name__ == "__main__":
    asyncio.run(main())
