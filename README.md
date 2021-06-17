![Catorce](docs/catorce.png)

---

Catorce is a telegram bot that can play UNO!

Inspired by [mau_mau_bot](https://github.com/jh0ker/mau_mau_bot), this was created to change some rules, avoid game resets as our group can take a very long time to play and add new features, like a scoreboard.

# Getting Started

To run this bot, make sure you have the latest version of go installed, than clone this repository, add a `.env` file to the root with `TELEGRAM_TOKEN` set (you can get one with [BotFather](https://core.telegram.org/bots)).

Features are still being implemented and lots of bugs are to be expected. The bot only responds in portuguese.

# Playing

## New Game

To create a new game, call `/new` in a group chat with the bot, players can then `/join` the game (for now, a player can't be in two games at the same time)

After all players have joined (need 2 to 10 players), just `/run` to start the game.

On your turn, type `@bot_user` and your cards should appear, selecting one should play it, selecting an invalid card (or any card out of your turn) will just show the current game status. There's also an icon to draw a card (more cards if the last card was a +4 or a chain of +2)

## Scores

> NYI
