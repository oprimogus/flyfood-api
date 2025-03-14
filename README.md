# Flyfood-api

Repositório da API do FlyFood

## Sobre o Projeto

A API do FlyFood é uma solução desenvolvida em Golang que serve como o backend para o aplicativo 
FlyFood. Este aplicativo visa facilitar a descoberta de restaurantes e lojas na região dos usuários. 
Ele fornece informações básicas sobre estabelecimentos, como disponibilidade, e permite que os usuários façam 
pedidos de delivery ou optem por ir ao local. Além disso, o aplicativo oferece uma funcionalidade de avaliação 
para que os consumidores possam compartilhar suas experiências.

Os principais atores do sistema são:

- **Consumidor:** O usuário final que utiliza o aplicativo para conhecer restaurantes e lojas, fazer 
- pedidos (delivery ou ir ao local) e obter informações como disponibilidade dos estabelecimentos.
  
- **Dono da Loja:** O proprietário do estabelecimento que utiliza a plataforma para receber pedidos, 
- fazer marketing e gerenciar informações do seu negócio.

- **Futuro Motoboy:** A API está projetada para, futuramente, incluir funcionalidades para motoboys, 
- que poderão receber informações sobre entregas a serem realizadas.

Esta API é a espinha dorsal da aplicação FlyFood, garantindo uma integração eficiente e escalável 
entre os diferentes componentes do sistema e seus usuários.

## Funcionalidades

- **Consulta de Estabelecimentos:** Permite aos consumidores buscar e visualizar informações sobre 
restaurantes e lojas na sua região.
- **Pedidos e Entregas:** Facilita a realização de pedidos de delivery e fornece informações 
sobre a disponibilidade dos estabelecimentos.
- **Avaliações:** Permite que os consumidores avaliem e comentem sobre suas experiências em 
diferentes locais.
- **Gerenciamento de Pedidos:** Oferece uma plataforma para donos de lojas 
gerenciarem pedidos e realizar atividades de marketing.


## Depêndencias

Obs: Configuradas como `tool` usando a nova feature do `go1.24.0`

### 1. Migrate CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Repositório com mais informações [aqui](https://github.com/golang-migrate/migrate/tree/v4.16.2/cmd/migrate)

### 2. SQLC

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Repositório com mais informações [aqui](https://docs.sqlc.dev/en/stable/overview/install.html)

### 3. swaggo

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Repositório com mais informações [aqui](https://github.com/swaggo/swag)

### 4. Air hot reload

```bash
go install github.com/air-verse/air@latest
```

Repositório com mais informações [aqui](https://github.com/cosmtrek/air)

### 5. validator

Documentação disponível [aqui](https://github.com/go-playground/validator)



## Primeira vez ao rodar o app localmente

1. Criar .env e preencher variáveis de ambiente

    ```bash
    cp .env.example .env
    ```

2. Instalar dependências do app

    ```bash
    make install
    ```

3. Subir banco de dados e demais containers

    ```bash
    make up
    ```

4. Executar migrations

    ```bash
    make migration-up
    ```

5. Mockar dados no banco local

    ```bash
    make mock-db
    ```

6. Rodar a API

    ```bash
    make dev
    ```

7. Acessar swagger [aqui](http://localhost:8080/api/v1/swagger/index.html#/)