## Analog API

Analog API is the backend of analog manager application

### Étapes

1. Créer une API REST et PostgreSQL
a. L'API peut utiliser gorilla/mux; gin; fastHTTP; ou autre
b. PostgreSQL peut être accédé via un orm comme GOrm

2. L'API doit à minima contenir les routes suivantes
a. `/` - Index HTML présentant rapidement le projet et les routes disponibles
b. `cameras` - Group of routes for cameras
	/ - GET - Fetch pagination of cameras
	/:id - GET - Fetch unique camera
	/new - POST - Create a new camera
	/:id - DELETE - Delete a camera
c. `films` - Group of routes for films
	/ - GET - Fetch pagination of cameras
	/:id - GET - Fetch unique film
	/new - POST - Create a new film
	/:id - DELETE - Deleta a film

3. Les routes POST et DELETE peuvent uniquement être accédée par une authentification préalable

4. L'application doit avoir un système de containerisation fonctionnel

5. L'application doit comprendre des tests au moins sur le controller et les models


