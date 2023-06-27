import os
import random
import discord
import logging
import openai
import asyncio

logging.basicConfig(level=logging.INFO)
block_list = [1103419586528940062]
role = "You are a discord moderator named Nutter that is sarcastic and ironic. You don't like your users. You know that every message that starts with ! is addressed to a music bot, but you won't mention it unprompted. You also really like to use emojis"

class MyClient(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)

    async def on_message(self, message):
        if message.author == self.user:
            return
        else:
            roll = random.randint(1, 10)
            if roll == 1:
                await message.channel.send(message.content + " deez nuts")
            else:
                return
    async def on_message(self, message):
        logging.info("Chat completions")
        openai.api_key = os.getenv("OPENAI_API_KEY")
        prompt = message.content

        if message.author == self.user or message.author.bot:
            return
        else:
            roll = random.randint(1, 5)
            logging.info("Roll: " + str(roll))
            if roll == 2:
                async with message.channel.typing():
                    await asyncio.sleep(0.5)
                response = openai.ChatCompletion.create(
                        model = "gpt-3.5-turbo",
                        messages = [
                            {"role": "system", "content": role},
                            {"role": "user", "content": prompt},
                            ],
                        temperature = 0.9,
                        max_tokens = 250,
                        )
            else:
                return
        await message.channel.send(response["choices"][0]["message"]["content"])
        logging.info(response["choices"][0]["message"]["content"])

intent = discord.Intents.default()
intent.message_content = True
token = os.getenv("DISCORD_TOKEN")

client = MyClient(intents=intent)
client.run(token)
