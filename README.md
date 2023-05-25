# GitLab Webhook to Feishu Bot

This project allows you to set up a webhook in GitLab that will send notifications to a Feishu bot when specific actions occur in your repository.

## Installation

1. Clone this repository in your local machine.
2. Create a new Feishu bot and copy the webhook url.
3. Set the environment variable FEISHU_BOT_ACCESS_TOKEN to the webhook url you copied.
4. Copy `.env-example` to `.env`
5. Open `.env` file and add your configurations.
6. Run `docker-compose up -d --build` to build the container and start the service
7. In your GitLab project, go to Settings > Integrations and enter the URL of your server followed by the endpoint https://YOUR_DOMAIN/gitlab/webhook (e.g. https://example.com/gitlab/webhook)
8. Select the events you want to trigger the webhook (e.g. Push events) and save the changes.


## Usage

When the events you selected are triggered, the webhook will be sent to your server and the server will be executed, sending a notification to your Feishu bot.

## TODO List
- [X] Push events
- [ ] Tag Push events
- [ ] Comments
- [ ] Confidential comments
- [ ] Issues events
- [ ] Confidential issues events
- [X] Merge request events
- [ ] Job events
- [ ] Pipeline events
- [ ] Wiki page events
- [ ] Deployment events
- [ ] Feature flag events
- [ ] Release events

## Contributing

Please feel free to contribute to this project by forking the repository and submitting a pull request. If you find any issues or have any suggestions for improvement, please open an issue. Thank you!

## License

This project is licensed under the MIT License.