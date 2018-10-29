
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

# Thoughts and execution of task
Uniqueness of timestamp = The timestamp in the application is currently not unique if two people upload in the same second.

Clock Trigger = Hosted on openstack and located in the *Clock Trigger* folder where it can be executed independent from the main API. I first developed it with direclty with
the datrabase, but changed it later to connect to heroku as this felt more meaningful to the total task.





