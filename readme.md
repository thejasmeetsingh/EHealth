# EHealth
It is an online booking management system where hospitals, clinics, or individual doctors can create an account and accept bookings, while users can view available health services near them and create bookings.

There is also a calendar management system where users can see the available slots of the doctors and create bookings as per their convenience, and the doctor calendar will be automatically blocked based on their individual bookings. Also, doctors can sync their Google/Outlook calendar events to block specific timings on their calendar.

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
- You can use this Postman Collection for using/testing the APIs: [Link](https://api.postman.com/collections/17396704-2cebc0d3-fcb4-4475-94b6-d2a59361463d?access_key=PMAT-01HDNGX7D80B3B81107SB6H9MZ)

## TODO
- [x] User management
- [ ] Booking management.
- [ ] Implement Auth system, JWT specifically.
- [ ] Google/Outlook calendar integration using OAuth.
- [ ] Implement geolocation so that users can view doctors that are nearest to them.
- [ ] Integrate an email service for sending booking-related notifications.