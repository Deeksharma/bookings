# Bookings and Reservations

This is a repository for bookings and registration project.

We have a bread and breakfast which has 2 rooms, and we want to make
it available online for people to book.

- Built in Go version 1.13
- uses [chi router](https://github.com/go-chi/chi)
- uses [alex edwards scs](https://github.com/alexedwards/scs) for session management
- uses [nosurf](https://github.com/justinas/nosurf) for CSRF validation
- uses [govalidator](https://github.com/asaskevich/govalidator) for email validation 
- uses [vanillajsdatepicker](https://mymth.github.io/vanillajs-datepicker) for picking dates
- uses [notie](https://github.com/jaredreich/notie) for sending alerts
- uses [sweetalert](https://sweetalert2.github.io/#input-types) for popup boxes

People can search for dates, look at what's available,
make a booking and reservation.

- website that showcases the property
- book a room, for more than one day - available or not available
- notify owner as well as guest
- for owner he'll login and sees whose coming and when they're coming
- review existing bookings, shows a calendar that'll show the bookings
- change or cancel the booking


We need :-
- Authentication system - only owner can look at bookings, change bookings or cancel them
- user can see their bookings and cancel as well
- store the info - database
- send notifications - email and text]

We are seeding the rooms table with data required using migration so that whenever we use command `soda reset`
we don't have to put the rooms again and again.

Send email to the owner and customer after reservation confirmation.

[foundation frontend framework](https://get.foundation/)
[foundation framework to format mails](https://get.foundation/emails)

// we'll create a new table that keeps the new reservations that are not processed yet, and that are shown to the owner when he clicks on new reservation - future reservations actually

[Data grid](https://github.com/fiduswriter/Simple-DataTables) that we are using

to get the project up:-
make database.yml, change the db name, username and password.
download mailhog using brew, brew services start mailhog, go to localhost:8025  