import os
import random
import discord

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

intent = discord.Intents.default()
intent.message_content = True
token = os.getenv("DISCORD_TOKEN")

client = MyClient(intents=intent)
client.run(token)
