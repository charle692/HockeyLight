# HockeyLight

This project basically recreates the Budweiser Red Light.

This project has 3 parts:
1. The Goal light itself
2. An app that is used to configure the goal light, which is found [here](https://github.com/charle692/HockeyLightMobile)
3. Auto wifi configuration using [pi-wifi](https://github.com/charle692/pi-wifi)

## Inspiration
- https://github.com/arim215/nhl_goal_light

## Installation
First off, this project is created for the Raspberry pi. However, it may work on other devices.

- Clone the repository
- Compile the code using [xgo](https://github.com/karalabe/xgo), `xgo -out main --targets=linux/arm "path to project root"`.
- Transfer the binary and the `hockey_light.db` to `/home/pi` on your Raspberry pi.
- Create a service that runs the binary
- Use the [mobile app](https://github.com/charle692/HockeyLightMobile) to configure the HockeyLight

## Hardware
- Raspberry Pi
- 12v Beacon
- Speaker
- 2 module relay
- 12v power adapter
- 12v -> 5v converter

## Algorithm
- Every day, check what games are being played. If the selected team is playing that day, record the time that the game will start.
- Once the game has started, every 3 seconds, pull from the api and check the scores.
- If the selected team has scored, trigger the goal light.
- At the end of the game, trigger the goal light if the selected team has won.
