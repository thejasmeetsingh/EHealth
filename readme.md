# EHealth
It is an online booking management system where hospitals, clinics, or individual doctors can create an account and accept bookings, while users can view available health services near them and create bookings.

## Follow below steps to setup and run the project:
- Install [docker](https://www.docker.com/products/docker-desktop/)
- Create a file named as `.env` in the project root path and add below varriables
    ```
    SECRET_KEY=

    PORT=

    DB_USER=
    DB_PASSWORD=
    DB_NAME=
    DB_URL=

    GIN_MODE=
    CGO_ENABLED=

    FROM_EMAIL=
    SMTP_SERVER=
    SMTP_PORT=
    SMTP_USERNAME=
    SMTP_PASSWORD=
    DEFAULT_RECIPIENT_EMAIL=

    ACCESS_TOKEN_EXP=
    REFRESH_TOKEN_EXP=
    ```
    You can set the respected values as per your need
- And finally run `docker-compose up` and that's it!

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/17396704-2cebc0d3-fcb4-4475-94b6-d2a59361463d?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D17396704-2cebc0d3-fcb4-4475-94b6-d2a59361463d%26entityType%3Dcollection%26workspaceId%3D392b781a-05ab-415b-9eb8-456aca6f3129)

## TODO
- [x] User management
- [x] Booking management.
- [x] Implement Auth system, JWT specifically.
- [x] Implement geolocation so that users can view doctors that are nearest to them.
- [x] Integrate an email service for sending booking-related notifications.