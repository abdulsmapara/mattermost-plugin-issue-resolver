<p align="center">
	<h1 align="center">Mattermost Issue Resolver Plugin</h1>
	<h5 align="center">A plugin that helps to resolve issues directly through Mattermost application</h5>
</p>


## Table of Content
- [About-the-plugin](#about-the-plugin)
- [Installation](#installation)
- [Working](#working)
	* [Managing Skills](#managing-skills)
	* [Resolving the issue](#resolving-the-issue)
- [Features](#features)
    * [Slash Commands](#slash-commands)
        + [```/skills```](#skills)
        + [```/resolve```](#resolve)
- [Hackathon-Mattermost-bot-hackfest](#hackathon-mattermost-bot-hackfest)
     * [Idea Phase Submission](#idea-phase-submission)
- [Additional Info](#additional-info)

## About the plugin
This plugin helps to resolve issues directly through Mattermost application. It enables the users to get solutions to their issues/questions if a similar/related question is found on [stackoverflow](https://www.stackoverflow.com). If a similar/related question is not found on stackoverflow, it suggests contacting a user who is skilled in the domain of the question/issue.

## Installation
Download the latest version of the [release](https://github.com/abdulsmapara/mattermost-plugin-issue-resolver/releases) directory. Go to `System Console` and upload the latest release in plugins section. For help on how to install a custom plugin, please refer [installing custom plugin docs](https://docs.mattermost.com/administration/plugins.html#custom-plugins).

*Currently unstable due to active development, should be used for testing purpose only*. 


## Working
1. This plugin works with the help of slash commands. Headover to any channel and write the slash command ```/resolve``` followed by your issue/question.
1. If a similar question/issue is found (that has an accepted answer) on stackoverflow, then the issue_resolver bot will reply with the question found and the accepted solution.
The following screenshot(s) show that the user's issue was "/resolve Java: Array out of bounds exception" which is searched on stackoverflow, and the complete question (along with its title) and the accepted solution on stackoverflow is fetched and displayed to the user by the bot. 
![User's issue](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot1.png)
![Similar/Related Question Found](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot2.png)
![Accepted Solution](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot3.png)

1. If a similar question/issue is not found or there is no accepted answer for the question on stackoverflow, then the bot tries to return a user who possesses the skill required to solve the issue.

The following screenshot shows that the user's issue was "call golang code from haskell", which we understand requires the skills "golang, haskell". The bot suggests a user who has mentioned "golang, haskell" among his/her skills.
![Suggestion-to-contact-a-User](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot4.png)

## Features
#### Slash Commands
```/skills```

The plugin suggests contacting a user who is skilled in the domain of the question/issue posted if a similar/related issue/question is not found on stackoverflow. In order to accomplish this task, users are required to update their skills.

In order to manage his/her skills, the user should use the slash command ```/skills``` and an optional command along with it. 

If no command is mentioned along with ```/skills```, the system lists the skills of the user stored.

If the command ```list <prefix (optional)>``` is used, the system prints the list of skills any user can have that starts with the given prefix. (If no prefix is mentioned, all the possible skills are listed). That is, there is a huge predefined set of skills, that any user can have. 

If the command ```add <comma separated list of skills>``` is used, the system adds all the skills for the user that are mentioned and are among the predefined list.

If the command ```delete <comma separated list of skills>``` is used, the system deletes all the skills mentioned for the user if he/she possessed that skill beforehand.

The GIF below depicts the usage of ```/skills```
![Usage of ```/skills```](https://github.com/abdulsmapara/Github-Media/blob/master/gif1.gif)

```/resolve <issue/question>```

The issue/question mentioned is searched on stackoverflow. If a similar question/issue is found that has an accepted answer on stackoverflow, then the issue_resolver bot will reply with the question found and the accepted solution.

If a similar question/issue with an accepted answer is not found, then a user will be searched who possesses the skill for resolving the issue. If such a user is found, the username of the user found is suggested by the bot for contact. 
The GIF below depicts the usage of ```/resolve```

![Usage of ```/resolve```](https://github.com/abdulsmapara/Github-Media/blob/master/gif2.gif)

## Hackathon [Mattermost Bot Hackfest](https://www.hackerearth.com/challenges/hackathon/mattermost-bot-hackfest/)
#### Idea Phase Submission

- Link to the ppt describing the idea - [Idea Presentation](https://he-s3.s3.amazonaws.com/media/sprint/mattermost-bot-hackfest/team/782765/8a5bcbfcodeblooded_mattermost_hackfest.pptx?Signature=YZdp812LgWXUaiup1j5GYe4TKQ8%3D&Expires=1583010208&AWSAccessKeyId=AKIA6I2ISGOYH7WWS3G5)




## Additional Info


This plugin is created for demonstration of the idea submitted at [Mattermost Bot Hackfest](#hackathon-mattermost-bot-hackfest).


The plugin is created by [@abdulsmapara](https://github.com/abdulsmapara).
