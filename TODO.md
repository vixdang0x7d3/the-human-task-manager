# short-term goal

Our short-term goal is to finish the personal task management module as fast as possible

pre-final test evaluation milestone: 2024/11/25

# to-do's

**Refactor**
- [x] Rename handlers, file, domain
- [x] Refactor routing, middleware, logging
- [x] Centralize server configuration logic
- [ ] Write better configuration logic 

**User domain**
- [x] impl user handler's create user http endpoint
- [x] impl user core's create user: hash password, call user store to save user to database
- [x] refactor unit tests
- [x] implement basic data validation (check empty fields) for user core's create user
- [x] refactor validation to app layer :D
- [x] impl user handlers' get user http endpoint, check if uuid is valid
- [x] impl user core's get user: call user store to get user via uuid.

**Task domain**
(nothing yet)

**UI**
- [x] setup templ
- [x] serve a basic hello world page
- [x] tailwindcss & theme


