# WP SSH/SFTP Backup CLI

With this CLI you will be able to connect through SSH to the desired client, generate a ZIP file in the SiteGround server of the "www/" directory, copy that file to the local directory under './backup', delete the server file, and download the zip through SFTP.

Then, upload this ZIP file to the specified GSC Bucket.

## TODO
- Test suite
- Upload the ZIP to the GCS Bucket
- Add an env file? so no credential input is needed
## Installation


```bash
go build main.go
```

## Example usage on the terminal

```bash
example:wp-backup-cli djigo$ go run main.go connect

Please write your SSH username:
sshUsername

Now introduce the client host:
partnerskiwi.com

Now introduce the desired port:
18765

Successfully connected to ssh server.
File will be called: 2021_03_04_160033_partnerskiwi.com.zip
Zipping the site, please be patient...
Done!
Copying file from SFTP to your ./backups directory ...
Also done!
Deleting zip file from the server...
Exiting, bye! :D
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
