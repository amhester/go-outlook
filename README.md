# go-outlook

An SDK for accessing Microsoft's graph API with Go.

I noticed while working on a project that there wasn't really any go sdk for outlook's api and so I decided to write this one. As of now, this package supports access to microsoft's graph api, exposing services for listing a user's calendars, their email folders, as well as exposing CRUD operations on calendar events and email messages. Also, as of right now, this project includes 0 dependencies.

Unfortunately, Microsoft's graph API exposes a lot more than what I currently support in this SDK, so this will very much be a work in progress. Also, for the time being, I have decided not to implement Microsoft's authentication flow in this package. However, the session service does require a user's refreshToken and will automatically handle fetching the access token for the given refreshToken.

## Installation

This is just a simple go package, so feel free to install via your tool of choice.

```bash
go get https://github.com/amhester/go-outlook
```

### Dependencies

None

### Environment Variables

This SDK requires the use of an authenticated session for all of it's exposed methods. Thus, it needs to be able to handle making requests on behalf of an application/user. To do that, the initial outlook client can be configured with both an App ID as well as an App Secret (provided by microsoft upon creation of an appliaction for their APIs). You can set these fields on the client either by passing them in as a ClientOpt on creation of the client, setting them after the client has been created, or through the following environment variables:

```bash
OUTLOOK_APP_ID=<YOUR_APPLICATION_ID>
OUTLOOK_APP_SECRET=<YOUR_APPLICATION_SECRET>
```

## Usage

Docs and Examples to come

## TODO

Write TODOs

## Testing

Yeah, probably still need to write some of those
