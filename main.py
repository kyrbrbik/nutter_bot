import os
import random
import discord
import openai
import logging

logging.basicConfig(level=logging.INFO)

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

        if message.author == self.user:
            return
        else:
            roll = random.randint(1, 5)
            if roll == 2:
                response = openai.ChatCompletion.create(
                        model = "gpt-3.5-turbo",
                        messages = [
                            {"role": "system", "content": "You are a discord moderator that sarcastically replies to user in his server"},
                            {"role": "user", "content": prompt},
                            ],
                        temperature = 0.9,
                        max_tokens = 250,
                        )
            else:
                return
        await message.channel.send(response["choices"][0]["message"]["content"])

intent = discord.Intents.default()
intent.message_content = True
token = os.getenv("DISCORD_TOKEN")

client = MyClient(intents=intent)
client.run(token)
