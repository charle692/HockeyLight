# HockeyLight

This project basically recreates the Budweiser Red Light.

This project has 2 parts:
1. The Goal light itself
2. An app that configured the goal light

## App Requirements
The app will need to be available for both Android and iOS. The focus at first is going to be Android. 
The app will be used to select which team the user is chearing for. The selected team is the only team that can trigger the goal light.

Additionaly the app will be used to configure the Raspberry pi's network settings a la chromecast. To achieve chromecast like configuratiion I'm most likely going to use https://github.com/jasbur/RaspiWiFi. 

Finally the app will allow the user to configure a delay. This delay is used to control when the light/horn goes off after a goal to prevent any spoilers. 

## Similar Project
- https://github.com/arim215/nhl_goal_light

## Hardware
- Raspberry Pi
- 12v Beacon
- Speaker
- 2 module relay
- 12v power adapter
- 12v -> 5v converter

## API
To get all the games for the specified day
https://statsapi.web.nhl.com/api/v1/schedule?startDate=2017-09-18&endDate=2017-09-18

To get data for the current game
https://statsapi.web.nhl.com/api/v1/game/2017010011/feed/live

## Algorithm
- Every day, check what games are being played. If the selected team is playing that day, record the time that the game will start.
- Once the game has started, every 3 seconds, pull from the api and check the scores.
- If the selected team has scored, trigger the goal light.
- At the end of the game, trigger the goal light if the selected team has won.
