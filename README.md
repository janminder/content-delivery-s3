# Static files from custom S3 Storage on Cloud Foundry

Basically this is a simple file proxy, which serves files from a third party s3 storage. Use case could be serving 
static files like javascript, css or other resources in a cloud foundry environment with 
a third-party or aws s3 storage layer.

```
              +-------+                 +--------------+
GET File      |       |                 |              |
+--------->   |  API  +<--------------->+  S3 Storage  |
              |       |                 |              |
              +-------+                 +---------+----+
                                                  ^
                                                  |
                                                  |
                                                  |
                                                  |
                                                  +
                                             Upload Files

```

### url schema for clients
https://static.domain.com/bucket-name/file-name

### Manage Data
Manage the data with a S3 Client:  https://cyberduck.io/download/

Cyberduck uses profile files to manage UI components. If you want to configure a third party S3 Storage (local or remote) you could load a generic profile. 

- [Generic S3 Profile for Cyberduck (http)](https://svn.cyberduck.io/trunk/profiles/S3%20(HTTP).cyberduckprofile)
- [Generic S3 Profile for Cyberduck (https)](https://svn.cyberduck.io/trunk/profiles/S3%20(HTTPS).cyberduckprofile)

### Build

Run `make` and a binary will be built to bin/cd-s3. Define your target arch in make file. 

### Deploy to Cloud Foundry

- a env variable named `profile` indicates the binary, the current environment
- the application detects based on this variable the cf env and search for an s3 service based on service name in `config.cloud.toml`

`cf push -f manifest.yml`

### Testing

This Project was tested in a cloud foundry environment with a Dell EMC Service.

### Work in progress..

- Manage the buckets, files and publish states from an admin UI
- Integrate the ability to use with a oAuth provider to manage admin permissions
