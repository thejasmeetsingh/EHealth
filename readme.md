# EHealth
It is an online booking management system where hospitals, clinics, or individual doctors can create an account and accept bookings, while users can view available health services near them and create bookings.

There is also a calendar management system where users can see the available slots of the doctors and create bookings as per their convenience, and the doctor calendar will be automatically blocked based on their individual bookings. Also, doctors can sync their Google/Outlook calendar events to block specific timings on their calendar.

## Follow below steps to setup and run the project:
- Install [docker](https://www.docker.com/products/docker-desktop/)
- Create a file named as `.env` in the project root path and add below varriables
    ```
    PORT=

    DB_USER=
    DB_PASSWORD=
    DB_NAME=
    DB_URL=

    GIN_MODE=
    CGO_ENABLED=
    ```
    You can set the respected values as per your need
- And finally run `docker-compose up` and that's it!

## TODO
- [ ] User management
- [ ] Booking management.
- [ ] Implement Auth system, JWT specifically.
- [ ] Google/Outlook calendar integration using OAuth.
- [ ] Implement geolocation so that users can view doctors that are nearest to them.
- [ ] Integrate an email service for sending booking-related notifications.