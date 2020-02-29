<p align="center">
	<h1 align="center">Mattermost Issue Resolver Plugin</h1>
	<h5 align="center">A plugin that helps resolve issues directly through Mattermost application</h5>
</p>


## Table of Content
- [About-the-plugin](#about-the-plugin)
- [Installation](#installation)
- [Prerequisites for getting the most of it](#prerequisites-for-getting-the-most-of-it)
- [Running up](#running-up)
## About the plugin
This plugin helps resolve issues directly through Mattermost application. It enables the users to get solutions to their issues/questions if found similar/related question is found on [stackoverflow](https://www.stackoverflow.com). If a similar/related question is not found on [stackoverflow](https://www.stackoverflow.com), it suggests contacting a user who is skilled in the domain of the question/issue.
## Installation
Download the latest version of the [release](https://github.com/abdulsmapara/mattermost-plugin-issue-resolver/releases) directory. Go to `System Console` and upload the latest release in plugins section. For help on how to install a custom plugin, please refer [installing custom plugin docs](https://docs.mattermost.com/administration/plugins.html#custom-plugins).

*Currently unstable due to active development, should be used for testing purpose only*. 

## Prerequisites for getting the most of it
1. The plugin suggests contacting a user who is skilled in the domain of the question/issue posted if a similar/related issue/question is not found on [stackoverflow](https://www.stackoverflow.com). 
1. To possess skills for the users, in order to accomplish the above mentioned task, users are required to update their skills.
1. In order to manage skills, user should use the slash command ```/skills``` and an optional command along with it.

## Running up
1. This plugin works with the help of simple commands. Headover to any channel and write the slash command ```/resolve``` followed by your issue/question, that will be searched on [stackoverflow](https://www.stackoverflow.com).