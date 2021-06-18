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

After all players have joined (need 2 to 10 players), just `/run` to start the game. For now, a player can't join a running game.

## Your Turn

On your turn, type `@bot_user` and your cards should appear, selecting one should play it, selecting an invalid card (or any card out of your turn) will just show the current game status. There's also an icon to draw a card (more cards if the last card was a +4 or a chain of +2)

## End

When someone plays their last card, the game finishes, points are calculated and statistics are updated. Type `/new` to start a new one.

### Points

Points are calculated according to the cards left on the hand when the game finishes:

| Card         | Value            |
| ------------ | ---------------- |
| Number Cards | Face Value (0-9) |
| Draw 2       | 20               |
| Reverse      | 20               |
| Skip         | 20               |
| Wild         | 50               |
| Draw Four    | 50               |

Less points = better

# Statistics

To see game statistics for the chat, use `/stats`, this will show some statistics like total games played and average response time, will also send a table with the chat ranking (based on average points).

You can use `/statsself` to see your own stats for the current chat
