import discord
import logging

from discord import commands

logging.basicConfig(level=logging.DEBUG)

client = commands.Bot(command_prefix = '.')

@client.event
async def on_ready()
    print('Bot is ready')

client.run('API_KEY')
