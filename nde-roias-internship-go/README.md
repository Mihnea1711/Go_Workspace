We utilize the public API available at https://dummyjson.com/products.
We retrieve the data from this API and populate our own database named project_publicapi.
The data is stored in the Product table.

# Project Structure:
- The "init_app" package serves to initialize the entire application. It verifies and initializes both the database and server. This package calls functions from the "database_controller" and "routes" packages.
- The "routes" package defines our routes along with the corresponding functions called for each route and method. These functions are invoked from the "controller" package.
- The "database_controller" package contains the methods for interacting with the database and its tables.
- The "jsonworking" package houses several functions for fetching information from Postman and sending responses back.
- The "filesUtils" folder contains JSON files that can be copied and pasted into Postman. These files serve as aids and are not directly utilized within the program.
- The "html" folder contains HTML files. Some functions are responsible for rendering HTML interfaces.
- The image "poza_API.jpg" illustrates the structure of the JSON provided by the public API.
