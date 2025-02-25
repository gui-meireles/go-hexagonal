# Arquitetura hexagonal

---

## Explicando a arquitetura:

### Pasta `application`

Vamos utilizar a pasta `application` para definir o **coração da nossa aplicação**, que ficará isolada do mundo externo.

Utilizaremos a `application/product.go` para definir as interfaces, estrutura do produto e as regras de negócio.

Em `application/product_service.go` criaremos os métodos para chamar a classe de persistência (que ficará fora da application) e conta
com algumas validações caso aconteça algum erro durante a persistência dos dados.

### Pasta `adapters`

Vamos utilizar a pasta `adapters` para definir os componentes externos da nossa aplicação, como a persistência no banco
de dados, filas, web, entre outros.

#### Pasta `adapters/cli`

Esse será o lado esquerdo da arquitetura hexagonal, que representará o cliente pelo CLI.

Utilizaremos a `adapters/cli/product.go` para mapear uma instrução de criar/buscar/atualizar um produto por linha de comando.

Em `adapters/cli/product_test.go`, criaremos os testes do adapter CLI.

> Podemos utilizar a biblioteca **Cobra** para **facilitar as chamadas CLI**, siga o módulo `Sobre o Cobra`
> que está logo abaixo.

#### Pasta `adapters/web`

Esse será o lado esquerdo da arquitetura hexagonal, que representará o cliente pela Web.

Em `adapters/web/server/server.go` criaremos as configurações de rotas e middlewares.

Em `adapters/web/handler/product.go` teremos os handlers, que funcionam semelhantes a uma controller, onde conseguimos
buscar ou criar um produto por http.

Em `adapters/web/handler/error_json.go` converteremos a mensagem de erro para o formato json.

Em `adapters/web/handler/error_json_test.go` testaremos a mensagem de erro em json.

> Para rodar o serviço web, utilizaremos o **Cobra** novamente, siga o módulo `Utilizando o Cobra para serviço http` que está logo abaixo.

Em `adapters/dto/product.go` criaremos o DTO da requisição `POST` para fazer o **bind** do JSON com o product da application.

#### Pasta `adapters/db`

Esse será o lado direito da arquitetura hexagonal, que fará a conexão com o db.

Utilizaremos a `adapters/db/product.go` para criar a persistência com o banco de dados.

O arquivo `sqlite.db` na raiz do projeto, **será o nosso banco de dados**, para **criar as tabelas** abra o bash e digite:
`sqlite3 sqlite.db` e rode os comando: `create table products(id string, name string, price float, status string);`

Em `adapters/db/product_test.go`, criaremos os testes no banco.

### Raiz do projeto

Na raiz do projeto contamos com o arquivo `mainteste.go` que é um modo de **criarmos uma conexão** dos `adapters` na `application`
e realizar uma inserção e atualização de um produto no banco.

Temos o arquivo `sqlite.db` que é utilizado na `mainteste.go` para salvar os **insert/update/delete** do banco.

Para rodar essa classe, basta digitar o comando no bash: `go run mainteste.go`.

---

### Como subir o projeto:

Abra um terminal na **pasta raiz**, e logo em seguida digite: `docker-compose up -d`.

Para checar se o docker está de pé, rode o comando: `docker ps`.

Acesse o bash com: `docker exec -it appproduct bash`.

**Para rodar as classes de testes Go**, rode o comando no mesmo bash: `go test ./...`.
> Caso queira rodar uma classe de teste em específico, utilize `go test ./application/{arquivo_teste}`.

![img.png](readme_images/img.png)

### Para evitar bugs de versionamento:

Tenha preferência por utilizar as **mesmas versões do arquivo** `go.mod` existente, **copie** todos os `require`,
apague o arquivo `go.mod` existente e cole no arquivo `go.mod` que você criará logo abaixo:

Vamos criar um arquivo **go.mod** dentro do bash com o comando: `go mod init github.com/gui-meireles/go-hexagonal`, para facilitar
o download dos pacotes que vamos trabalhar.

---

### Comandos úteis do application

Para criar mocks facilmente em Go, podemos utilizar o **_mockgen_**, que é um auto gerador de mocks, para utilizá-lo
digite o comando a seguir: `mockgen -destination=application/mocks/application.go -source=application/product.go application`
dentro do bash do container.

### Comandos úteis do db

**Para criar o arquivo sqlite.db**, utilize o comando dentro do bash: `touch sqlite.db`.

Para abrir o terminal do db, rode o comando dentro do bash: `sqlite3 sqlite.db`.

Para checar as tabelas: `.tables`.

---

### Sobre o Cobra

É uma **biblioteca popular em Go** que permite criar aplicativos de linha de comando de forma _fácil e eficiente_.

Fornecendo uma estrutura para definir comandos, flags e argumentos.

Para utilizá-lo em sua aplicação Go, você deve baixá-lo através do **Dockerfile**. (Utilize a versão v1.1.3)

E inicie o Cobra com o comando: `cobra init --pkg-name=github.com/gui-meireles/go-hexagonal`.

Ele criará os arquivos `main.go` e `cmd/root.go`.

Caso você tente rodar o `go run main.go` e apareça **erros de pacotes**, rode o comando: `go mod tidy`, ele removerá
os pacotes que não estiverem sendo utilizados e baixará os necessários.

Ao rodar o comando: `go run main.go`, deverá aparecer semelhante a imagem abaixo:
![img_1.png](readme_images/img_1.png)

Crie o arquivo cli com: `cobra add cli` e ele ficará dentro da pasta `cmd`.

Faremos a injeção das dependências da service e do db no arquivo `cmd/root.go` e criaremos os comandos em
`cmd/cli.go`.

Com isso feito, podemos **ver os comandos CLI** com: `go run main.go cli --help`.

E para **executar uma ação de criar produto**, você pode rodar o comando: `go run main.go cli -a=create -n="Melancia" -p=25.0`.

| Comandos | Função                                             |
|----------|----------------------------------------------------|
| **-a**   | Chama a função dentro de `adapters/cli/product.go` |
| **-n**   | Informa o nome do produto                          |
| **-p**   | Informa o preço do produto                         |
| **-i**   | Informa o ID do produto                            |

Fazendo o **select** na tabela, podemos ver que foi criado um novo produto:
![img_2.png](readme_images/img_2.png)

Buscando pelo **ID do produto** através do **CLI**:
![img_3.png](readme_images/img_3.png)

### Utilizando o Cobra para serviço http

O Cobra pode nos auxiliar a criar um serviço http, para isso precisamos digitar no bash: `cobra add http`.

E dentro da pasta `cmd` será criado um arquivo `http.go`, onde instanciaremos nossa web service.

Com o arquivo `http.go` configurado, digite no bash: `go run main.go http` para iniciar o server.
> **Obs:** O comando acima precisa estar rodando no bash para que sejam feitas as requisições http.
![img_4.png](readme_images/img_4.png)

Em outro terminal, você pode chamar a requisição `GET` com o comando: `curl http://localhost:9000/product/{productId}`.
![img_5.png](readme_images/img_5.png)

Podemos utilizar o **Postman** para fazer requisições `POST` para criar um produto:
![img_6.png](readme_images/img_6.png)

Caso seja feita uma requisição em que o **price é negativo**, ele retornará **erro**:
![img_7.png](readme_images/img_7.png)

Para habilitar um produto existente, utilize: `http://localhost:9000/product/{productId}/enable`:
![img_8.png](readme_images/img_8.png)

E para desabilitar um produto existente, utilize: `http://localhost:9000/product/{productId}/disable`:
> Porém devemos zerar o valor do produto para desabilitá-lo.
![img_9.png](readme_images/img_9.png)