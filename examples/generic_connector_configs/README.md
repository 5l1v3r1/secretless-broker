# Generic HTTP Connector Example Configurations

## Table of Contents
- [Generic HTTP Connector Example Configurations](#generic-http-connector-example-configurations)
  * [Table of Contents](#table-of-contents)
- [Sample Configurations](#sample-configurations)
  * [***OAuth 2.0 API***](#oauth-20-api)
      - [How to use this connector](#how-to-use-this-connector)
      - [Example Usage](#example-usage)
  * [***Github API***](#github-api)
      - [How to use this connector](#how-to-use-this-connector-1)
      - [Example Usage](#example-usage-1)
  * [***Elasticsearch API***](#elasticsearch-api)
      - [How to use this connector](#how-to-use-this-connector-2)
      - [Example Usage](#example-usage-2)
  * [***Slack Web API***](#slack-web-api)
      - [How to use this connector](#how-to-use-this-connector-3)
      - [Example Usage](#example-usage-3)
  * [***Splunk API***](#splunk-api)
      - [How to use this connector](#how-to-use-this-connector-4)
      - [Example Usage](#example-usage-4)
- [Contributing](#contributing)

# Sample Configurations

The [generic HTTP connector](../../internal/plugin/connectors/http/generic/README.md)
enables using Secretless with a wide array of HTTP-based services _without
having to write new Secretless connectors_. Instead, you can modify your
Secretless configuration to specify the header structure the HTTP service
requires to authenticate.

If your target uses self-signed certs you will need to follow the
[documented instructions](https://docs.secretless.io/Latest/en/Content/References/connectors/scl_handlers-https.htm#Manageservercertificates) for adding the target’s CA to
Secretless’ trusted certificate pool.

> Note: The following examples use the [Keychain provider](https://docs.cyberark.com/Product-Doc/OnlineHelp/AAM-DAP/11.3/en/Content/References/providers/scl_keychain.htm?TocPath=Fundamentals%7CSecretless%20Pattern%7CSecret%20Providers%7C_____5).
> Replace the service prefix `service#` with an appropriate service
> or use a different provider as needed.
___

## OAuth 2.0 API
This generic OAuth HTTP connector can be used for any service that accepts a
Bearer token as an authorization header.

The configuration file for the OAuth 2.0 API can be found at
[oauth2_secretless.yml](./oauth2_secretless.yml).

#### How to use this connector
* Edit the supplied service configuration to get your OAuth token
* Run secretless with the supplied configuration(s)
* Query the API using `http_proxy=localhost:8071 curl <Your OAuth2 API Endpoint URL>/{Request}`

#### Example Usage

___
## Github API
This example can be used to interact with [Github's API](https://developer.github.com/v3/).

The configuration file for the Github API can be found at [github_secretless.yml](./github_secretless.yml).

#### How to use this connector
* Edit the supplied configuration to get your [GitHub OAuth token](https://developer.github.com/v3/#oauth2-token-sent-in-a-header) from the correct provider/path.
* Run Secretless with the supplied configuration
* Query the GitHub API using `http_proxy=localhost:8081 curl api.github.com/{request}`

#### Example Usage

___
## Elasticsearch API
This example can be used to interact with [Elasticsearch's API](https://www.elastic.co/guide/en/elasticsearch/reference/current).

The configuration file for the Elasticsearch API can be found at
[elasticsearch_secretless.yml](./elasticsearch_secretless.yml).

#### How to use this connector
* Edit the supplied configuration to get your Elasticsearch [API Key](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-create-api-key.html) or
[OAuth token](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-get-token.html)
* Run Secretless with the supplied configuration(s)
* Query the Elasticsearch API using `http_proxy=localhost:9020 curl <Elasticsearch Endpoint URL>/{Request}`

#### Example Usage
___

## Slack Web API
This example can be used to interact with [Slack's Web API](https://api.slack.com/apis).

The configuration file for the Slack Web API can be found at [slack_secretless.yml](./slack_secretless.yml).

#### How to use this connector
* Edit the supplied configuration to get your Slack [OAuth token](https://api.slack.com/legacy/oauth#flow)
* Run secretless with the supplied configuration(s)
* Query the Slack API using `http_proxy=localhost:9030 curl -d {data} <Slack Endpoint URL>` or `http_proxy=localhost:9040 curl -d {data} <Slack Endpoint URL>` depending on if your endpoint requires JSON or URL encoded requests

#### Example Usage

___
## Splunk API
This example can be used to interact with [Splunk's API](https://api.slack.com/apis).

The configuration file for the Slack Web API can be found at [splunk_secretless.yml](./splunk_secretless.yml).

#### How to use this connector

#### Example Usage
* Run a local instance of Splunk in a Docker container
```
docker run \
  -d \
  -p 8000:8000 \
  -p 8089:8089 \
  -e "SPLUNK_START_ARGS=--accept-license" \
  -e "SPLUNK_PASSWORD=specialpass" \
  --name splunk \
  splunk/splunk:latest
```
* Follow the instructions here to create a local Splunk token using Splunk Web.
* Save the local token from Splunk Web into the OSX keychain
* Add 'SplunkServerDefaultCert' at IP 127.0.0.1 to etc/hosts on the machine. This was so the host name of the HTTP Request would match the name on the certificate that is provided on our Splunk container.
* Usethe provided cacert.pem file on the Splunk docker container for my certificate, and write it to the local machine
```
docker exec -it splunk sudo cat /opt/splunk/etc/auth/cacert.pem > myLocalSplunkCertificate.pem
```
* Set a variable in the terminal named `SECRETLESS_HTTP_CA_BUNDLE` and set it to the path where myLocalSpunkCertificate.pem was on the local machine.
* Run Secretless
```
./dist/darwin/amd64/secretless-broker \
  -f examples/generic_connector_configs/splunk_secretless.yml
```
* On another terminal window, make a request to Splunk using Secretless
```
http_proxy=localhost:8081 curl -k -X GET SplunkServerDefaultCert:8089/services/apps/local
```


___

# Contributing

Do you have an HTTP service that you use? Can you write a Secretless generic
connector config for it? **Add the sample config to this folder and list it in
the table above!** Others may find your connector config useful, too - [send us
a PR](https://github.com/cyberark/community/blob/master/CONTRIBUTING.md#contribution-workflow)!
