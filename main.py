import os
import random
import discord
import logging
import openai
import asyncio

logging.basicConfig(level=logging.INFO)
block_list = [1103419586528940062]
role = "You are a discord moderator named Nutter that is sarcastic and ironic. You don't like your users. You also really like to use emojis"
is_waiting = False

class MyClient(discord.Client):

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
        global is_waiting
        if is_waiting == True:
            logging.info("Already waiting for response")
            return
        else:
            prompt = message.content
            if message.author == self.user or message.author.bot or message.content.startswith("!"):
                logging.info("Message from self or bot")
                return
            else:
                roll = self.dice_roll()
                logging.info(roll)
                if roll == 2:
                    response = self.api_call(prompt)
                    logging.info("Response: " + response["choices"][0]["message"]["content"])
                    await message.channel.send(response["choices"][0]["message"]["content"])
                else:
                    return
    def api_call(self, prompt):
        global is_waiting
        is_waiting = True
        prompt = prompt
        logging.info(prompt)
        openai.api_key = os.getenv("OPENAI_API_KEY")
        response = openai.ChatCompletion.create(
            model = "gpt-3.5-turbo",
            messages = [
                {"role": "system", "content": role},
                {"role": "user", "content": prompt},
                ],
            temperature = 0.9,
            max_tokens = 250,
            )
        is_waiting = False
        return response
    def dice_roll(self):
        roll = random.randint(1, 2)
        return roll


intent = discord.Intents.default()
intent.message_content = True
intent.guilds = True
token = os.getenv("DISCORD_TOKEN")

client = MyClient(intents=intent)

@client.event
async def on_ready():
    logging.info("Logged in as {0.user}".format(client))
    guilds = client.guilds
    for guild in guilds:
        logging.info(guild.name)

client.run(token)


