- We have the public API https://dummyjson.com/products
- We retrieve the data from it and populate our own database named project_publicapi
- This data goes into the Product table.


Project structure:

- We have the "init_app" package that initializes the entire application. It checks and initializes the database and server. It calls functions from the "database_controller" and "routes" packages.
- We have the "routes" package where our routes are defined along with the functions that are called for each route and method. The functions from the "controller" package are called from here.
- We have the "database_controller" package where the methods for the database and tables are defined.
- Avem package-ul "jsonworking" unde sunt cateva functii pentru a lua informatii din postman si pentru a trimite un raspuns in postman
- in folder-ul "filesUtils" sunt json-uri pe care le dam copy paste in postman. (e ca un ajutor. Nu se folosesc in program filele respective)
- In the "html" folder, there are HTML files. Some functions are used to render HTML interfaces.
- The image "poza_API.jpg" presents how the JSON offered by the public API looks like.




