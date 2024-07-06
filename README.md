# MS PORTFOLIO BS

## Instalación y ejecución ⚙️

- Crear archivo .env usando .env.template

  ```bash
  #COMMON
  APP_NAME=ms-portfolio-bs
  HTTP_PORT=3002

  #MONGODB
  MONGO_INITDB_ROOT_USERNAME=admin
  MONGO_INITDB_ROOT_PASSWORD=123456
  MONGO_INITDB_DATABASE=portfolio_db

  MONGO_URI="mongodb://admin:123456@localhost:27018/portfolio_db?authSource=admin"
  ```

- Levantar proyecto

  ```bash
  air
  ```

- [Ir a la documentación en Swagger](http://localhost:3002/swagger/index.html)

  ![swagger](/etc/images/swagger.png)
  
## Endpoints 🚀

- Client Portfolio Detail
  ![portfolio_detail](/etc/images/portfolio_detail.png)
  
- Seed
  ![seed](/etc/images/seed.png)

## Construido con 🛠️

- Go
- Gin
- MongoDB
- Swagger
- Docker

### Autor ✒️

- Williams David Galeano Gomez, <willyrhcp96@gmail.com>
