# short-term goal

Our short-term goal is to finish the personal task management module as fast as possible

pre-final test evaluation milestone: 2024/11/25

# to-do's

**User domain**
- [x] impl user handler's create user http endpoint
- [x] impl user core's create user: hash password, call user store to save user to database
- [x] refactor unit tests
- [x] implement basic data validation (check empty fields) for user core's create user
- [ ] refactor validation to app layer :D
- [ ] impl user handlers' get user http endpoint, check if uuid is valid
- [ ] impl user core's get user: call user store to get user via uuid.
- [ ] wire up the endpoints in main


**Task domain**
(nothing yet)

**UI**
- [ ] setup templ and tailwindcss
- [ ] serve a basic hello world page

