<p align="center">
	<h1 align="center">Mattermost Issue Resolver Plugin</h1>
	<h5 align="center">A plugin that helps resolve issues directly through Mattermost application</h5>
</p>


## Table of Content
- [About-the-plugin](#about-the-plugin)
- [Installation](#installation)
- [Working](#working)
	* [Managing Skills](#managing-skills)
	* [Resolving the issue](#resolving-the-issue)
## About the plugin
This plugin helps resolve issues directly through Mattermost application. It enables the users to get solutions to their issues/questions if found similar/related question is found on [stackoverflow](https://www.stackoverflow.com). If a similar/related question is not found on [stackoverflow](https://www.stackoverflow.com), it suggests contacting a user who is skilled in the domain of the question/issue.
## Installation
Download the latest version of the [release](https://github.com/abdulsmapara/mattermost-plugin-issue-resolver/releases) directory. Go to `System Console` and upload the latest release in plugins section. For help on how to install a custom plugin, please refer [installing custom plugin docs](https://docs.mattermost.com/administration/plugins.html#custom-plugins).

*Currently unstable due to active development, should be used for testing purpose only*. 


## Working
1. This plugin works with the help of slash commands. Headover to any channel and write the slash command ```/resolve``` followed by your issue/question.
1. If a similar question/issue is found (that has an accepted answer) on stackoverflow, then the issue_resolver bot will reply with the question found and the accepted solution.
The following screenshot(s) show that the user's issue was "Java: Array Index Out of Bounds" which is searched on stackoverflow, and the complete question (along with its title) and the accepted solution on stackoverflow is fetched and displayed to the user by the bot. 
![User's issue](https://drive.google.com/uc?export=view&id=1cjntMynrKxl-RP5U3sHL1d4HvrwZqO0I)
![Similar/Related Question Found](https://drive.google.com/uc?export=view&id=1sM0QVFEN87oV5rek5sNozSrgybuurjDq)
![Accepted Solution](https://drive.google.com/uc?export=view&id=1Fh4Yali4fe49dzvu5iZWK9qZFXSRQtyP)

1. If a similar question/issue is not found or there is no accepted answer for the question on stackoverflow, then the bot tries to return a user who possesses the skill required to solve the issue.

The following screenshot shows that the user's issue was "call c function from app python", which we understand requires the skills "C, Python". The bot suggests a user who has mentioned "C, Python" among his/her skills.
![Suggestion-to-contact-a-User](https://drive.google.com/uc?export=view&id=1vZJFZpR3SA6C2pznnfq73-WCdCYh8K5Z)

## Features
#### Slash Commands
1. ```/skills```

The plugin suggests contacting a user who is skilled in the domain of the question/issue posted if a similar/related issue/question is not found on [stackoverflow](https://www.stackoverflow.com). In order to accomplish this task, users are required to update their skills.
In order to manage his/her skills, the user should use the slash command ```/skills``` and an optional command along with it. 
	2. If no command is mentioned along with ```/skills```, the system lists the skills of the user stored.
	2. If the command ```list <prefix (optional)>``` is used, the system prints the list of skills any user can have that starts with the given prefix. (If no prefix is mentioned, all the skills are listed). That is, there is a huge predefined set of skills, that any user can have. 
	2. If the command ```add <comma separated list of skills>``` is used, the system adds all the skills for the user that are mentioned and are among the predefined list.
	2. If the command ```delete <comma separated list of skills>``` is used, the system deletes all the skills mentioned for the user if he/she possessed that skill beforehand.

1. ```/resolve <issue/question>```

Issue/Question will be searched on stackoverflow. If a similar question/issue is found that has an accepted answer on stackoverflow, then the issue_resolver bot will reply with the question found and the accepted solution.
If a similar question/issue with an accepted answer is not found, then a user will be searched who possesses the skill for resolving the issue. If such a user is found, the username of the user found is suggested by the bot for contact. 
