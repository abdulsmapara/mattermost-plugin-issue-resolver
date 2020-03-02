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
        + ```/skills```
        + ```/resolve```
- [Development](#development)
- [Hackathon-Mattermost-bot-hackfest](#hackathon-mattermost-bot-hackfest)
     * [Idea Phase Submission](#idea-phase-submission)
     * [Project Submission](#project-submission)
- [Demo Video](#demo-video)
- [Additional Info](#additional-info)

## About the plugin
This plugin helps to resolve issues directly through Mattermost application. It enables the users to get solutions to their issues/questions if a similar/related question is found on [stackoverflow](https://www.stackoverflow.com). If a similar/related question is not found on stackoverflow, it suggests contacting a user who is skilled in the domain of the question/issue.

## Installation
Download the latest version of the [release](https://github.com/abdulsmapara/mattermost-plugin-issue-resolver/releases) directory. Go to `System Console` and upload the latest release in plugins section. For help on how to install a custom plugin, please refer [installing custom plugin docs](https://docs.mattermost.com/administration/plugins.html#custom-plugins).

## Working
1. This plugin works with the help of slash commands. Headover to any channel and write the slash command ```/resolve``` followed by your issue/question.
1. If a similar question/issue is found (that has an accepted answer) on stackoverflow, then the issue_resolver bot will reply with the question found and the accepted solution.
The following screenshot(s) show that the user's issue was "/resolve Java: Array out of bounds exception" which is searched on stackoverflow, and the complete question (along with its title) and the accepted solution on stackoverflow is fetched and displayed to the user by the bot. 
![User's issue](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot1.png)
![Similar/Related Question Found](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot2.png)
![Accepted Solution](https://github.com/abdulsmapara/Github-Media/blob/master/screenshot3.png)

1. If a similar question/issue is not found by the `issue_resolver` bot or there is no accepted answer for the question on stackoverflow, then the bot tries to return a user who possesses the skill required to solve the issue.

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

The issue/question mentioned is searched on stackoverflow. If a similar question/issue is found that has an accepted answer on stackoverflow, then the `issue_resolver` bot will reply with the question found and the accepted solution.

If a similar question/issue with an accepted answer is not found, then a user will be searched who possesses the skill for resolving the issue. If such a user is found, the username of the user found is suggested by the bot for contact. 
The GIF below depicts the usage of ```/resolve```

![Usage of ```/resolve```](https://github.com/abdulsmapara/Github-Media/blob/master/gif2.gif)

#### Note : 
1. The plugin logs important events at appropriate log levels.
1. The plugin is tested on [Mattermost Server version 5.19.0](https://github.com/mattermost/mattermost-server/releases/tag/v5.19.0)

## Development

1. This plugin contains only the server.

1. Use `make check-style` to check the style.

1. Use `make test` to test the plugin.

1. Use `make dist` to build distributions of the plugin that can be uploaded to a Mattermost server

1. Alternatively, use `make` to check the style, test and build distributions of the plugin that you can upload to a Mattermost server (all at once).

1. Use `make deploy` to deploy the plugin to your local server. Before running `make deploy`, you need to set a few environment variables:
	```
		export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
		export MM_ADMIN_USERNAME=admin
		export MM_ADMIN_PASSWORD=password
	```

1. If you want to deploy the plugin by using `System Console`:
	
	1. On the the server, in the file `config/config.json`, change `EnableUploads` in the `Plugin Settings` to `true`

	1. Login to Mattermost server with admin privileges.

	1. Headover to `System Console` and upload the tar.gz file created in `dist/` directory to the plugins section.
	For help on how to install a custom plugin, please refer [installing custom plugin docs](https://docs.mattermost.com/administration/plugins.html#custom-plugins).

	1. Enable the plugin in the section `Installed Plugins` on the same page.


## Hackathon [Mattermost Bot Hackfest](https://www.hackerearth.com/challenges/hackathon/mattermost-bot-hackfest/)
#### Idea Phase Submission

- Link to the ppt describing the idea - [Idea Presentation](https://he-s3.s3.amazonaws.com/media/sprint/mattermost-bot-hackfest/team/782765/8a5bcbfcodeblooded_mattermost_hackfest.pptx?Signature=QVDAn%2F2SPbab8UJbYXv2Jb%2FWjBA%3D&Expires=1583084663&AWSAccessKeyId=AKIA6I2ISGOYH7WWS3G5)

#### Project Submission

- Link to the submitted file - [File for project submission](https://github.com/mattermost/mattermost-hackathon-hackerearth-jan2020/blob/master/hackathon-submissions/abdulsmapara-mattermost-plugin-issue-resolver.md)

- Link to the pull request for submitting the project - [Project Submission](https://github.com/mattermost/mattermost-hackathon-hackerearth-jan2020/pull/13)

## Demo Video

Link - [https://vimeo.com/394781065/](https://vimeo.com/394781065/)

## Additional Info


This plugin is created for demonstration of the idea submitted at [Mattermost Bot Hackfest](#hackathon-mattermost-bot-hackfest).


The plugin is created by [@abdulsmapara](https://github.com/abdulsmapara).
