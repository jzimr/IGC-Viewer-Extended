
![Scheme](logo.png)

IGC Viewer Extended is an online service that allows users to browse information about IGC files.

# Author
Jan Zimmer (janzim, 493594)

# Links
Heroku app: [http://igcviewer-extended.herokuapp.com/paragliding/api](http://igcviewer-extended.herokuapp.com/paragliding/api)  
Discord Webhook: [https://discordapp.com/api/webhooks/505722994237374466/6yqNRGY1b8jitN_jyhHxLhGc-xThQBqW3L0-xC-X86Hrd__Zi_eMAGki87lv5xzbY2IQ](https://discordapp.com/api/webhooks/505722994237374466/6yqNRGY1b8jitN_jyhHxLhGc-xThQBqW3L0-xC-X86Hrd__Zi_eMAGki87lv5xzbY2IQ)  
Discord Channel (For Clock Trigger): [https://discord.gg/n9T2Qjr](https://discord.gg/n9T2Qjr)

# Configuration of API
Both the *main* API and the *Clock Trigger* have a *config.json* configuration file where the administrator can set the database information, how many tracks are shown  
per page, the root for admin API and webhook URL.

# Choices and decisions
**Uniqueness of timestamp** = The timestamp in the application is currently not unique if two people upload in the same second. However, unless we want to get a specific timestamp,
this does not matter since they are all listed in an array upon calling the API endpoints.

**Clock Trigger** = Hosted on openstack and located in the *Clock Trigger* folder where it can be executed independently from the main API. At first it communicated directly with
the tracks in the database, but I changed this later on to connect to heroku as this felt more meaningful overall. 
The clock trigger contains a *config.json* where the admin can change the webhook URL, as this makes it much easier to setup rather than hardcoding it into the source code.

**Choice of webhooks** = I've chosen to only go with the Discord webhook as it is easier to check for messages and people can easily check the results by clicking the link for the
Discord channel above.

**Paging cap** = The paging cap (default=5) can be modified in the *config.json* by changing the value.





