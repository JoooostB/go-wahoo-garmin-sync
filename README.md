# Syncs files from Wahoo/Dropbox to Garmin
This project lets you synchronise your activities from Wahoo to Garmin Connect using Dropbox. Whenever a new activity gets uploaded in Dropbox by the ELMNT app, the webhook will be called which downloads the file from Dropbox and uploads it to Garmin Connect using [`abrander/garmin-connect`](https://github.com/abrander/garmin-connect).

## Configuration
1. Configure your ELMNT app to push your activities to Dropbox.
2. Create a Developer app on Dropbox
3. Add app key and secret in .envrc
4. Start application 
5. Visit application to create oauth2 session
6. Upload new activity to Dropbox

## Development
Fill the neccessary variables in .envrc and run the application with Air and Redis using docker-compose up.