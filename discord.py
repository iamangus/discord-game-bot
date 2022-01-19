import discord
import logging
import argparse

from discord.ext import commands
logging.basicConfig(level=logging.DEBUG)

# Initialize parser
parser = argparse.ArgumentParser()
# Adding optional argument
parser.add_argument("-t", "--bottoken", help = "Discord bot token")
# Read arguments from command line
args = parser.parse_args()

client = commands.Bot(command_prefix = '.')

@client.event
async def on_ready():
    print('Bot is ready.')

client.run('args.bottoken')

 

